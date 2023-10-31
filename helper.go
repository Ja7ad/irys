package irys

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/signer"
	"github.com/Ja7ad/irys/types"
)

func decodeBody[T any](body io.Reader) (T, error) {
	var resp T
	d := json.NewDecoder(body)
	return resp, d.Decode(&resp)
}

func addContentType(contentType string, tags ...types.Tag) types.Tags {
	found := false
	for _, tag := range tags {
		if tag.Name == "Content-Type" {
			found = true
		}
	}

	if !found {
		tags = append(tags, types.Tag{Name: "Content-Type", Value: contentType})
	}

	return tags
}

func statusCheck(resp *http.Response) error {
	switch {
	case resp.StatusCode == http.StatusPaymentRequired:
		return errors.ErrNotEnoughBalance
	case resp.StatusCode >= http.StatusBadRequest:
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%d: %s", resp.StatusCode, string(b))
	case resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted:
		return nil
	}
	return nil
}

func signFile(file []byte, signer signer.Signer, withAnchor bool, tags ...types.Tag) ([]byte, error) {
	tags = addContentType(http.DetectContentType(file), tags...)

	dataItem := types.BundleItem{
		Data: types.Base64String(file),
		Tags: tags,
	}

	if withAnchor {
		anchor := make([]byte, 32)
		_, err := rand.Read(anchor)
		if err != nil {
			return nil, err
		}
		dataItem.Anchor = anchor
	}

	dataItem.Anchor.Base64()

	if err := dataItem.Sign(signer); err != nil {
		return nil, err
	}

	reader, err := dataItem.Reader()
	if err != nil {
		return nil, err
	}

	signedByte, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return signedByte, nil
}
