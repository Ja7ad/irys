package currency

import (
	"os"

	"github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/signer"
)

const (
	_arweave_name   = "arweave"
	_arweave_chain  = "arweave"
	_arweave_symbol = "ar"
)

type Arweave struct {
	unimplementedEther
	chain     string
	symbol    string
	name      string
	rpc       string
	tokenType CurrencyType
	signer    *signer.ArweaveSigner
}

// NewArweaveFromFile create token object for arweave by private key file arweave (not supported for TopUp Balance)
func NewArweaveFromFile(filePath, rpc string) (Currency, error) {
	privateKey, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	s, err := signer.NewArweaveSigner(string(privateKey))
	if err != nil {
		return nil, err
	}

	return &Arweave{
		chain:     _arweave_chain,
		symbol:    _arweave_symbol,
		name:      _arweave_name,
		tokenType: ARWEAVE,
		rpc:       rpc,
		signer:    s,
	}, nil
}

// NewArweave create token object from arweave private key payload  (not supported for TopUp Balance)
func NewArweave(privateKey string) (Currency, error) {
	if len(privateKey) == 0 {
		return nil, errors.ErrPrivateKeyIsEmpty
	}

	s, err := signer.NewArweaveSigner(privateKey)
	if err != nil {
		return nil, err
	}

	return &Arweave{
		chain:  _arweave_chain,
		symbol: _arweave_symbol,
		name:   _arweave_name,
		signer: s,
	}, nil
}

func (a *Arweave) GetChain() string {
	return a.chain
}

func (a *Arweave) GetSymbol() string {
	return a.symbol
}

func (a *Arweave) GetName() string {
	return a.name
}

func (a *Arweave) GetSinger() signer.Signer {
	return a.signer
}

func (a *Arweave) GetRPCAddr() string {
	return a.rpc
}

func (a *Arweave) GetType() CurrencyType {
	return a.tokenType
}
