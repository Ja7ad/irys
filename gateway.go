package irys

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/Ja7ad/irys/errors"

	"github.com/warp-contracts/syncer/src/utils/arweave"
	"github.com/warp-contracts/syncer/src/utils/bundlr"
)

const (
	_pricePath    = "%s/price/%s/%v"
	_uploadPath   = "%s/tx/%s"
	_txPath       = "%s/tx/%s"
	_downloadPath = "%s/%s"
)

func (i *Irys) GetPrice(fileSize int) (*big.Int, error) {
	url := fmt.Sprintf(_pricePath, i.network, i.token.GetBundlrName(), fileSize)
	resp, err := i.client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(b))
	}

	return decodeBody[*big.Int](resp.Body)
}

func (i *Irys) Upload(file io.Reader, tags ...Tag) (Transaction, error) {
	url := fmt.Sprintf(_uploadPath, i.network, i.token.GetBundlrName())
	b, err := io.ReadAll(file)
	if err != nil {
		return Transaction{}, err
	}
	return i.upload(url, b, tags...)
}

func (i *Irys) Download(hash string) (*File, error) {
	url := fmt.Sprintf(_downloadPath, i.network, hash)
	resp, err := i.client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(b))
	}

	return &File{
		Data:          resp.Body,
		Header:        resp.Header,
		ContentLength: resp.ContentLength,
		ContentType:   resp.Header.Get("Content-Type"),
	}, nil
}

func (i *Irys) GetMetaData(hash string) (Transaction, error) {
	url := fmt.Sprintf(_txPath, i.network, hash)
	resp, err := i.client.Get(url)
	if err != nil {
		return Transaction{}, err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return Transaction{}, err
		}
		return Transaction{}, fmt.Errorf("%d: %s", resp.StatusCode, string(b))
	}

	return decodeBody[Transaction](resp.Body)
}

func (i *Irys) upload(url string, payload []byte, tags ...Tag) (Transaction, error) {
	var signer bundlr.Signer
	bundlrTags := make(bundlr.Tags, 0)

	for _, tag := range tags {
		bundlrTags = append(bundlrTags, bundlr.Tag(tag))
	}

	bundlrTags = addContentType(http.DetectContentType(payload), bundlrTags)

	switch i.token.GetBundlrName() {
	case "matic", "ethereum":
		s, err := bundlr.NewEthereumSigner(i.token.GetPrivate())
		if err != nil {
			return Transaction{}, err
		}
		signer = s
	case "arweave":
		s, err := bundlr.NewArweaveSigner(i.token.GetPrivate())
		if err != nil {
			return Transaction{}, err
		}
		signer = s
	default:
		return Transaction{}, errors.ErrBundlrIsInvalid
	}

	dataItem := bundlr.BundleItem{
		Data: arweave.Base64String(payload),
		Tags: bundlrTags,
	}

	if err := dataItem.Sign(signer); err != nil {
		return Transaction{}, err
	}

	reader, err := dataItem.Reader()
	if err != nil {
		return Transaction{}, err
	}

	b, err := io.ReadAll(reader)
	if err != nil {
		return Transaction{}, err
	}

	body := bytes.NewBuffer(b)

	resp, err := i.client.Post(url, "application/octet-stream", body)
	if err != nil {
		return Transaction{}, err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return Transaction{}, err
		}
		return Transaction{}, fmt.Errorf("%d: %s", resp.StatusCode, string(b))
	}

	return decodeBody[Transaction](resp.Body)
}

func decodeBody[T any](body io.Reader) (T, error) {
	var resp T
	d := json.NewDecoder(body)
	return resp, d.Decode(&resp)
}

func addContentType(contentType string, tags bundlr.Tags) []bundlr.Tag {
	found := false
	for _, tag := range tags {
		if tag.Name == "Content-Type" {
			found = true
		}
	}

	if !found {
		tags = append(tags, bundlr.Tag{Name: "Content-Type", Value: contentType})
	}

	return tags
}
