package token

const _0x_prefix = "0x"

type Token interface {
	GetBundlrName() string
	GetChain() string
	GetSymbol() string
	GetPrivate() string
}
