package errors

import "errors"

var (
	ErrTestNetRPCNotSet  = errors.New("testnet rpc has been not set, please set using option TestNetRPC")
	ErrPrivateKeyIsEmpty = errors.New("private key is empty")
	ErrTokenNotSupported = errors.New("token not supported in this method")
	ErrBundlrIsInvalid   = errors.New("bundlr link is invalid")
)
