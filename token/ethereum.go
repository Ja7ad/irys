package token

import "github.com/Ja7ad/irys/errors"

const (
	_ethereum_name   = "ethereum"
	_ethereum_chain  = "ethereum"
	_ethereum_symbol = "eth"
)

type Ethereum struct {
	chain      string
	symbol     string
	name       string
	privateKey string
}

// NewEthereum create ethereum token object
func NewEthereum(privateKey string) (Token, error) {
	if len(privateKey) == 0 {
		return nil, errors.ErrPrivateKeyIsEmpty
	}

	return &Ethereum{
		chain:      _ethereum_chain,
		symbol:     _ethereum_symbol,
		name:       _ethereum_name,
		privateKey: _0x_prefix + privateKey,
	}, nil
}

func (e *Ethereum) GetChain() string {
	return e.chain
}

func (e *Ethereum) GetSymbol() string {
	return e.symbol
}

func (e *Ethereum) GetBundlrName() string {
	return e.name
}

func (e *Ethereum) GetPrivate() string {
	return e.privateKey
}
