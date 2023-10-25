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

	_fantom_name   = "fantom"
	_fantom_chain  = "fantom"
	_fantom_symbol = "fantom"

	_matic_name   = "matic"
	_matic_chain  = "polygon"
	_matic_symbol = "matic"

	_bnb_name   = "bnb"
	_bnb_chain  = "binance"
	_bnb_symbol = "bnb"

	_arbitrum_name   = "arbitrum"
	_arbitrum_chain  = "arbitrum"
	_arbitrum_symbol = "arbitrum"

	_avalanche_name   = "avalanche"
	_avalanche_chain  = "avalanche"
	_avalanche_symbol = "avalanche"
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

	return &Ethereum{
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

// NewBNB create bnb object token
func NewBNB(privateKey, rpc string) (Token, error) {
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
		chain:      _bnb_chain,
		symbol:     _bnb_symbol,
		name:       _bnb_name,
		signer:     s,
		tokenType:  MATIC,
		rpc:        rpc,
		client:     client,
		privateKey: prKey,
		publicKey:  publicKeyECDSA,
	}, nil
}

// NewFantom create fantom object token
func NewFantom(privateKey, rpc string) (Token, error) {
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
		chain:      _fantom_chain,
		symbol:     _fantom_symbol,
		name:       _fantom_name,
		signer:     s,
		tokenType:  FANTOM,
		rpc:        rpc,
		client:     client,
		privateKey: prKey,
		publicKey:  publicKeyECDSA,
	}, nil
}

// NewAvalanche create avalanche object token
func NewAvalanche(privateKey, rpc string) (Token, error) {
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
		chain:      _avalanche_chain,
		symbol:     _avalanche_symbol,
		name:       _avalanche_name,
		signer:     s,
		tokenType:  AVALANCHE,
		rpc:        rpc,
		client:     client,
		privateKey: prKey,
		publicKey:  publicKeyECDSA,
	}, nil
}

// NewArbitrum create arbitrum object token
func NewArbitrum(privateKey, rpc string) (Token, error) {
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
		chain:      _arbitrum_chain,
		symbol:     _arbitrum_symbol,
		name:       _arbitrum_name,
		signer:     s,
		tokenType:  ARBITRUM,
		rpc:        rpc,
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
