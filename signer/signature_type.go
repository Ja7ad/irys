package signer

import "strconv"

type SignatureType int

// Values are taken from bundlr library
// https://github.com/Bundlr-Network/arbundles/blob/5413fe576098355f7502a5fa9456f8db6a861492/src/constants.ts#L4
const (
	Arweave SignatureType = iota + 1
	APTOS
	Ethereum
	SOLANA
	NEAR
)

func (self SignatureType) Bytes() []byte {
	return []byte(strconv.Itoa(int(self)))
}
