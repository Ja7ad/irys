package irys

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	errs "github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/types"
	"github.com/avast/retry-go"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
)

func (c *Client) ChunkUpload(ctx context.Context, file io.Reader, tags ...types.Tag) (types.Transaction, error) {
	var (
		wg   sync.WaitGroup
		once sync.Once
	)
	jobsCh := make(chan types.Job)
	errCh := make(chan error)
	workerNum := _defaultNumWorkers

	b, err := signFile(file, c.currency.GetSinger(), tags...)
	if err != nil {
		return types.Transaction{}, err
	}

	chunkSize := len(b) / workerNum

	switch {
	case len(b) < _defaultMinChunk || chunkSize > _defaultMaxChunk:
		return types.Transaction{}, errs.ErrNotAllowedChunkSize
	case len(b) > _defaultMinChunk && len(b) < 1000000:
		chunkSize = len(b)
		workerNum = 1
	case chunkSize < _defaultMinChunk:
		return types.Transaction{}, errs.ErrNotAllowedChunkSize
	}

	chunkInfo, err := generateChunkUUID(ctx, c)
	if err != nil {
		return types.Transaction{}, err
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

	for start := 0; start < len(b); start += chunkSize {
		end := start + chunkSize
		if end > len(b) {
			end = len(b)
		}

		chunkData := b[start:end]
		chunk := types.Chunk{ID: chunkInfo.ID, Offset: int64(index * chunkSize), Data: chunkData}
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
		if err := finishChunk(ctx, c, chunkInfo.ID); err != nil {
			return types.Transaction{}, err
		}

		return getChunkTx(ctx, c, chunkInfo.ID)
	}
}

func generateChunkUUID(ctx context.Context, c *Client) (types.ChunkResponse, error) {
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
	data, err := json.Marshal(chunk)
	if err != nil {
		return err
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
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

func finishChunk(ctx context.Context, c *Client, uuid string) error {
	url := fmt.Sprintf(_chunkUpload, c.network, c.currency.GetName(), uuid, -1)

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, url, nil)
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
		c.debugMsg("[ChunkUpload] finalizing chunk upload")
		return statusCheck(resp)
	}
}

func getChunkTx(ctx context.Context, c *Client, uuid string) (types.Transaction, error) {
	url := fmt.Sprintf(_chunkUpload, c.network, c.currency.GetName(), uuid, 0)

	var tx types.Transaction
	err := retry.Do(func() error {
		req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("x-chunking-version", "2")

		resp, err := c.client.Do(req)
		if err != nil {
			return err
		}

		if err := statusCheck(resp); err != nil {
			return err
		}

		transaction, err := decodeBody[types.Transaction](resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		tx = transaction
		return nil
	}, retry.Attempts(3))

	if err != nil {
		return types.Transaction{}, err
	}

	return tx, nil
}
