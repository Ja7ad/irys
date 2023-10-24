package irys

import (
	"github.com/Ja7ad/irys/token"
	"io"
	"math/big"
	"net/http"
)

type Irys struct {
	client  *http.Client
	network Network
	token   token.Token
}

type Gateway interface {
	// GetPrice return fee base on fileSize in byte for selected token
	GetPrice(fileSize int) (*big.Int, error)
	// Upload upload file
	Upload(file io.Reader, tags ...Tag) (Transaction, error)
	// Download get file with header details
	Download(hash string) (*File, error)
	// GetMetaData get transaction details
	GetMetaData(hash string) (Transaction, error)
}

// New create Irys object
func New(network Network, token token.Token, options ...Option) (Gateway, error) {
	irys := new(Irys)
	irys.client = http.DefaultClient
	irys.network = network
	irys.token = token

	for _, opt := range options {
		opt(irys)
	}

	return irys, nil
}
