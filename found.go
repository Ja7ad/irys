package irys

import (
	"context"
	"github.com/Ja7ad/irys/currency"
	"github.com/Ja7ad/irys/errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"math/big"
)

func (i *IrysClient) createTx(ctx context.Context, amount *big.Int) (string, error) {
	switch i.currency.GetType() {
	case currency.ETHEREUM, currency.MATIC, currency.AVALANCHE, currency.FANTOM, currency.BNB, currency.ARBITRUM:
		hash, err := createEthTx(ctx, i, amount)
		if err != nil {
			return "", err
		}
		return hash, nil
	//TODO: arweave not supported currently
	case currency.ARWEAVE:

	}
	return "", errors.ErrTokenNotSupported
}

func createEthTx(ctx context.Context, i *IrysClient, amount *big.Int) (string, error) {
	pubKey := i.currency.GetPublicKey()
	client := i.currency.GetRPCClient()
	fromAddress := crypto.PubkeyToAddress(*pubKey)
	toAddress := common.HexToAddress(i.contract)

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return "", err
	}

	var data []byte
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	data = append(data, methodID...)
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	data = append(data, paddedAddress...)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	data = append(data, paddedAmount...)

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})

	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(i.contract),
		amount,
		gasLimit,
		gasPrice,
		data,
	)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), i.currency.GetPrivateKey())
	if err != nil {
		return "", err
	}

	if err = client.SendTransaction(ctx, signedTx); err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}
