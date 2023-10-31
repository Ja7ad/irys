package irys

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	errs "github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/types"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
)

const (
	_maxRetries      = 3 // define the maximum number of retries for a timeout error
	_defaultMinChunk = 500000
	_defaultMaxChunk = 95000000
)

func (c *Client) ChunkUpload(ctx context.Context, file io.Reader, chunkId string, tags ...types.Tag) (types.Transaction, error) {
	var (
		wg   sync.WaitGroup
		once sync.Once
	)
	jobsCh := make(chan types.Job)
	errCh := make(chan error)
	workerNum := 1
	chunkSize := 0
	chunkUUID := chunkId

	payload, err := io.ReadAll(file)
	if err != nil {
		return types.Transaction{}, err
	}

	b, err := signFile(payload, c.currency.GetSinger(), true, tags...)
	if err != nil {
		return types.Transaction{}, err
	}

	fileSize := len(b)

	if fileSize < _defaultMinChunk {
		return types.Transaction{}, errs.ErrNotAllowedChunkSize
	}

	switch {
	case fileSize >= 1000000 && fileSize < 10000000:
		workerNum = 2
	case fileSize >= 10000000 && fileSize < _defaultMaxChunk:
		workerNum = 3
	case fileSize >= _defaultMaxChunk:
		workerNum = 5
	}

	chunkSize = fileSize / workerNum

	if len(chunkUUID) == 0 {
		chunkInfo, err := generateChunkID(ctx, c)
		if err != nil {
			return types.Transaction{}, err
		}
		chunkUUID = chunkInfo.ID
	} else {
		// TODO: implement exists chunkId for resume
	}

	for w := 0; w < workerNum; w++ {
		go func(workerId int) {
			c.debugMsg("[ChunkUpload] create worker %v", workerId)
			if err := worker(ctx, c, workerId, &wg, jobsCh); err != nil {
				errCh <- err
			}
		}(w)
	}

	index := 0

	for start := 0; start < fileSize; start += chunkSize {
		end := start + chunkSize
		if end > fileSize {
			end = fileSize
		}

		chunkData := b[start:end]
		chunk := types.Chunk{ID: chunkUUID, Offset: int64(index * chunkSize), Data: chunkData}
		job := types.Job{Chunk: chunk, Index: index}
		jobsCh <- job
		index++
		wg.Add(1)
		c.debugMsg("[ChunkUpload] create job with index %v", index)
	}

	go func() {
		wg.Wait()
		close(jobsCh)
		close(errCh)
	}()

	for e := range errCh {
		once.Do(func() {
			close(jobsCh)
			close(errCh)
		})
		return types.Transaction{}, e
	}

	select {
	case <-ctx.Done():
		return types.Transaction{}, ctx.Err()
	default:
		return finishChunk(ctx, c, chunkUUID)
	}
}

func generateChunkID(ctx context.Context, c *Client) (types.ChunkResponse, error) {
	url := fmt.Sprintf(_chunkUpload, c.network, c.currency.GetName(), -1, -1)

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return types.ChunkResponse{}, err
	}

	req.Header.Set("x-chunking-version", "2")

	resp, err := c.client.Do(req)
	if err != nil {
		return types.ChunkResponse{}, err
	}
	defer resp.Body.Close()

	if err := statusCheck(resp); err != nil {
		return types.ChunkResponse{}, err
	}

	return decodeBody[types.ChunkResponse](resp.Body)
}

func getChunkID(ctx context.Context, c *Client, chunkId string) (types.ChunkInfoResponse, error) {
	panic("implement me")
}

func worker(ctx context.Context, c *Client, id int, wg *sync.WaitGroup, jobs <-chan types.Job) error {
	defer wg.Done()
	for job := range jobs {
		numTries := 0
		for numTries < _maxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				err := createChunkRequest(ctx, c, job.Chunk, job.Index, id)
				// if we have a network timeout error, retry the request
				var urlErr *url.Error
				if errors.As(err, &urlErr) {
					var netErr net.Error
					if errors.As(urlErr.Err, &netErr) && netErr.Timeout() {
						numTries++
						c.debugMsg("[ChunkUpload] timeout occurred during execution chunk upload, retrying... (Attempt %d of %d)", numTries, _maxRetries)
						continue
					}
				} else {
					return err
				}
				break
			}
		}
	}
	return nil
}

func createChunkRequest(ctx context.Context, c *Client, chunk types.Chunk, index, workerID int) error {
	url := fmt.Sprintf(_chunkUpload, c.network, c.currency.GetName(), chunk.ID, chunk.Offset)

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(chunk.Data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("x-chunking-version", "2")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.debugMsg("[ChunkUpload] worker %d do request for chunk %d", workerID, index)
		return statusCheck(resp)
	}
}

func finishChunk(ctx context.Context, c *Client, uuid string) (types.Transaction, error) {
	url := fmt.Sprintf(_chunkUpload, c.network, c.currency.GetName(), uuid, -1)

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return types.Transaction{}, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("x-chunking-version", "2")

	resp, err := c.client.Do(req)
	if err != nil {
		return types.Transaction{}, err
	}

	defer resp.Body.Close()

	select {
	case <-ctx.Done():
		return types.Transaction{}, ctx.Err()
	default:
		if err := statusCheck(resp); err != nil {
			return types.Transaction{}, err
		}

		return decodeBody[types.Transaction](resp.Body)
	}
}
