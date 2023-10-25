package irys

import (
	"encoding/json"
	"github.com/Ja7ad/irys/types"
	"io"
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
