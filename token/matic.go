package token

import "github.com/Ja7ad/irys/errors"

const (
	_matic_name   = "matic"
	_matic_chain  = "ethereum"
	_matic_symbol = "matic"
)

type Matic struct {
	chain      string
	symbol     string
	name       string
	privateKey string
}

// NewMatic create matic object token
func NewMatic(privateKey string) (Token, error) {
	if len(privateKey) == 0 {
		return nil, errors.ErrPrivateKeyIsEmpty
	}

	return &Matic{
		chain:      _matic_chain,
		symbol:     _matic_symbol,
		name:       _matic_name,
		privateKey: _0x_prefix + privateKey,
	}, nil
}

func (m *Matic) GetChain() string {
	return m.chain
}

func (m *Matic) GetSymbol() string {
	return m.symbol
}

func (m *Matic) GetBundlrName() string {
	return m.name
}

func (m *Matic) GetPrivate() string {
	return m.privateKey
}
