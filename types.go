package irys

import (
	"io"
	"net/http"
)

type Transaction struct {
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Address   string `json:"address"`
	Owner     string `json:"owner"`
	Signature string `json:"signature"`
	Target    string `json:"target"`
	Tags      []Tag  `json:"tags"`
	Anchor    string `json:"anchor"`
	DataSize  string `json:"data_size"`
	RawSize   string `json:"raw_size"`
}

type Tag struct {
	Name  string `json:"name" avro:"name"`
	Value string `json:"value" avro:"value"`
}

type File struct {
	Data          io.ReadCloser
	Header        http.Header
	ContentLength int64
	ContentType   string
}
