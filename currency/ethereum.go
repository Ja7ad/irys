package currency

import (
	"crypto/ecdsa"

	"github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/signer"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Ethereum struct {
	chain      string
	symbol     string
	name       string
	rpc        string
	tokenType  CurrencyType
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	signer     *signer.EthereumSigner
}

// NewEthereum create ethereum currency object
func NewEthereum(privateKey, rpc string) (Currency, error) {
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
		name:       "ethereum",
		chain:      "ethereum",
		symbol:     "eth",
		signer:     s,
		rpc:        rpc,
		tokenType:  ETHEREUM,
		client:     client,
		privateKey: prKey,
		publicKey:  publicKeyECDSA,
	}, nil
}

// NewMatic create matic object currency
func NewMatic(privateKey, rpc string) (Currency, error) {
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
		name:       "matic",
		chain:      "polygon",
		symbol:     "matic",
		signer:     s,
		tokenType:  MATIC,
		rpc:        rpc,
		client:     client,
		privateKey: prKey,
		publicKey:  publicKeyECDSA,
	}, nil
}

// NewBNB create bnb object currency
func NewBNB(privateKey, rpc string) (Currency, error) {
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
		name:       "bnb",
		chain:      "binance",
		symbol:     "bnb",
		signer:     s,
		tokenType:  BNB,
		rpc:        rpc,
		client:     client,
		privateKey: prKey,
		publicKey:  publicKeyECDSA,
	}, nil
}

// NewArbitrum create arbitrum object currency
func NewArbitrum(privateKey, rpc string) (Currency, error) {
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
		name:       "arbitrum",
		chain:      "arbitrum",
		symbol:     "arb",
		signer:     s,
		tokenType:  ARBITRUM,
		rpc:        rpc,
		client:     client,
		privateKey: prKey,
		publicKey:  publicKeyECDSA,
	}, nil
}

// NewAvalanche create avalanche object currency
func NewAvalanche(privateKey, rpc string) (Currency, error) {
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
		name:       "avalanche",
		chain:      "avalanche",
		symbol:     "avax",
		signer:     s,
		tokenType:  AVALANCHE,
		rpc:        rpc,
		client:     client,
		privateKey: prKey,
		publicKey:  publicKeyECDSA,
	}, nil
}

// NewFantom create fantom object currency
func NewFantom(privateKey, rpc string) (Currency, error) {
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
		name:       "fantom",
		chain:      "fantom",
		symbol:     "ftm",
		signer:     s,
		tokenType:  FANTOM,
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

func (e *Ethereum) GetName() string {
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

func (e *Ethereum) GetType() CurrencyType {
	return e.tokenType
}
