package token

import (
	"os"

	"github.com/Ja7ad/irys/errors"
)

const (
	_arweave_name   = "arweave"
	_arweave_chain  = "arweave"
	_arweave_symbol = "ar"
)

type Arweave struct {
	chain      string
	symbol     string
	name       string
	privateKey string
}

// NewArweaveFromFile create token object for arweave by private key file arweave
func NewArweaveFromFile(filePath string) (Token, error) {
	privateKey, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return &Arweave{
		chain:      _arweave_chain,
		symbol:     _arweave_symbol,
		name:       _arweave_name,
		privateKey: string(privateKey),
	}, nil
}

// NewArweave create token object from arweave private key payload
func NewArweave(privateKey string) (Token, error) {
	if len(privateKey) == 0 {
		return nil, errors.ErrPrivateKeyIsEmpty
	}

	return &Arweave{
		chain:      _arweave_chain,
		symbol:     _arweave_symbol,
		name:       _arweave_name,
		privateKey: privateKey,
	}, nil
}

func (a *Arweave) GetChain() string {
	return a.chain
}

func (a *Arweave) GetSymbol() string {
	return a.symbol
}

func (a *Arweave) GetBundlrName() string {
	return a.name
}

func (a *Arweave) GetPrivate() string {
	return a.privateKey
}
