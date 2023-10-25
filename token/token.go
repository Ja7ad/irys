package token

import (
	"crypto/ecdsa"
	"github.com/Ja7ad/irys/signer"
	"github.com/ethereum/go-ethereum/ethclient"
)

const _0x_prefix = "0x"

type TokenType uint8

const (
	ETHEREUM TokenType = iota
	MATIC
	BNB
	ARWEAVE
)

type Token interface {
	Ether
	GetBundlrName() string
	GetChain() string
	GetSymbol() string
	GetSinger() signer.Signer
	GetRPCAddr() string
	GetType() TokenType
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
