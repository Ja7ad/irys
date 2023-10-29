package irys

import (
	"encoding/json"
	"fmt"
	"github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/signer"
	"github.com/Ja7ad/irys/types"
	"io"
	"net/http"
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

func signFile(file io.Reader, signer signer.Signer, tags ...types.Tag) ([]byte, error) {
	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	tags = addContentType(http.DetectContentType(b), tags...)

	dataItem := types.BundleItem{
		Data: types.Base64String(b),
		Tags: tags,
	}

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
