package currency

import (
	"crypto/ecdsa"
	"github.com/Ja7ad/irys/signer"
	"github.com/ethereum/go-ethereum/ethclient"
)

const _0x_prefix = "0x"

type CurrencyType uint8

const (
	ETHEREUM CurrencyType = iota
	MATIC
	BNB
	ARBITRUM
	AVALANCHE
	FANTOM
	ARWEAVE
)

type Currency interface {
	Ether
	GetName() string
	GetChain() string
	GetSymbol() string
	GetSinger() signer.Signer
	GetRPCAddr() string
	GetType() CurrencyType
}

type Ether interface {
	GetRPCClient() *ethclient.Client
	GetPrivateKey() *ecdsa.PrivateKey
	GetPublicKey() *ecdsa.PublicKey
}

type unimplementedEther struct{}

var _ Ether = (*unimplementedEther)(nil)

func (u unimplementedEther) GetRPCClient() *ethclient.Client {
	panic("implement me")
}

func (u unimplementedEther) GetPrivateKey() *ecdsa.PrivateKey {
	panic("implement me")
}

func (u unimplementedEther) GetPublicKey() *ecdsa.PublicKey {
	panic("implement me")
}
