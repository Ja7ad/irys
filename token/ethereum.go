package token

import (
	"crypto/ecdsa"
	"github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/signer"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	_ethereum_name   = "ethereum"
	_ethereum_chain  = "ethereum"
	_ethereum_symbol = "eth"
)

type Ethereum struct {
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

// NewEthereum create ethereum token object
func NewEthereum(privateKey, rpc string) (Token, error) {
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

	return &Ethereum{
		chain:      _ethereum_chain,
		symbol:     _ethereum_symbol,
		name:       _ethereum_name,
		signer:     s,
		rpc:        rpc,
		tokenType:  ETHEREUM,
		client:     client,
		privateKey: prKey,
		publicKey:  publicKeyECDSA,
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

func (e *Ethereum) GetSinger() signer.Signer {
	return e.signer
}

func (e *Ethereum) GetRPCAddr() string {
	return e.rpc
}

func (e *Ethereum) GetRPCClient() *ethclient.Client {
	return e.client
}

func (e *Ethereum) GetPrivateKey() *ecdsa.PrivateKey {
	return e.privateKey
}

func (e *Ethereum) GetPublicKey() *ecdsa.PublicKey {
	return e.publicKey
}

func (e *Ethereum) GetType() TokenType {
	return e.tokenType
}
