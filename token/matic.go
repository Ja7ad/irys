package token

import (
	"crypto/ecdsa"
	"github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/signer"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	_matic_name   = "matic"
	_matic_chain  = "ethereum"
	_matic_symbol = "matic"
)

type Matic struct {
	chain      string
	symbol     string
	name       string
	rpc        string
	tokenType  TokenType
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	signer     *signer.EthereumSigner
}

// NewMatic create matic object token
func NewMatic(privateKey, rpc string) (Token, error) {
	if len(privateKey) == 0 {
		return nil, errors.ErrPrivateKeyIsEmpty
	}

	s, err := signer.NewEthereumSigner(_0x_prefix + privateKey)
	if err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	prKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	pbKey := prKey.Public()
	publicKeyECDSA, ok := pbKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.ErrAssertionPublicKey
	}

	return &Matic{
		chain:      _matic_chain,
		symbol:     _matic_symbol,
		name:       _matic_name,
		signer:     s,
		tokenType:  MATIC,
		rpc:        rpc,
		client:     client,
		privateKey: prKey,
		publicKey:  publicKeyECDSA,
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

func (m *Matic) GetSinger() signer.Signer {
	return m.signer
}

func (m *Matic) GetRPCAddr() string {
	return m.rpc
}

func (m *Matic) GetRPCClient() *ethclient.Client {
	return m.client
}

func (m *Matic) GetPrivateKey() *ecdsa.PrivateKey {
	return m.privateKey
}

func (m *Matic) GetPublicKey() *ecdsa.PublicKey {
	return m.publicKey
}

func (m *Matic) GetType() TokenType {
	return m.tokenType
}
