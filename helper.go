package irys

import (
	"encoding/json"
	"fmt"
	"github.com/Ja7ad/irys/errors"
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
