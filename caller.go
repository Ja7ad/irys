package irys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/types"
	retry "github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"math/big"
	"net/http"
)

const (
	_pricePath       = "%s/price/%s/%v"
	_uploadPath      = "%s/tx/%s"
	_txPath          = "%s/tx/%s"
	_downloadPath    = "%s/%s"
	_sendTxToBalance = "%s/account/balance/matic"
	_getBalance      = "%s/account/balance/matic?address=%s"
)

func (i *IrysClient) GetPrice(ctx context.Context, fileSize int) (*big.Int, error) {
	url := fmt.Sprintf(_pricePath, i.network, i.currency.GetName(), fileSize)

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if err := statusCheck(resp); err != nil {
			return nil, err
		}

		return decodeBody[*big.Int](resp.Body)
	}
}

func (i *IrysClient) GetBalance(ctx context.Context) (*big.Int, error) {
	pbKey := i.currency.GetPublicKey()
	url := fmt.Sprintf(_getBalance, i.network, crypto.PubkeyToAddress(*pbKey).Hex())

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if err := statusCheck(resp); err != nil {
			return nil, err
		}

		b, err := decodeBody[types.BalanceResponse](resp.Body)
		if err != nil {
			return nil, err
		}

		return b.ToBigInt(), nil
	}
}

func (i *IrysClient) TopUpBalance(ctx context.Context, amount *big.Int) (types.TopUpConfirmation, error) {
	urlConfirm := fmt.Sprintf(_sendTxToBalance, i.network)

	hash, err := i.createTx(ctx, amount)
	if err != nil {
		return types.TopUpConfirmation{}, err
	}

	b, err := json.Marshal(&types.TxToBalanceRequest{
		TxId: hash,
	})
	if err != nil {
		return types.TopUpConfirmation{}, err
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, urlConfirm, bytes.NewBuffer(b))
	if err != nil {
		return types.TopUpConfirmation{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := i.client.Do(req)
	if err != nil {
		return types.TopUpConfirmation{}, err
	}

	defer resp.Body.Close()

	select {
	case <-ctx.Done():
		return types.TopUpConfirmation{}, ctx.Err()
	default:
		if err := statusCheck(resp); err != nil {
			return types.TopUpConfirmation{}, err
		}

		confirm, err := decodeBody[types.TopUpConfirmationResponse](resp.Body)
		if err != nil {
			return types.TopUpConfirmation{}, err
		}

		if confirm.Confirmed {
			var balance *big.Int
			err = retry.Do(func() error {
				balance, err = i.GetBalance(ctx)
				return err
			}, retry.Attempts(3))

			if err != nil {
				return types.TopUpConfirmation{}, err
			}

			return types.TopUpConfirmation{
				Confirmed: true,
				Hash:      hash,
				Balance:   balance,
			}, nil
		}

		return types.TopUpConfirmation{}, nil
	}
}

func (i *IrysClient) BasicUpload(ctx context.Context, file io.Reader, tags ...types.Tag) (types.Transaction, error) {
	url := fmt.Sprintf(_uploadPath, i.network, i.currency.GetName())
	b, err := io.ReadAll(file)
	if err != nil {
		return types.Transaction{}, err
	}

	price, err := i.GetPrice(ctx, len(b))
	if err != nil {
		return types.Transaction{}, err
	}

	balance, err := i.GetBalance(ctx)
	if err != nil {
		return types.Transaction{}, err
	}

	if balance.Cmp(price) < 0 {
		_, err := i.TopUpBalance(ctx, price)
		if err != nil {
			return types.Transaction{}, err
		}
	}

	return i.upload(ctx, url, b, tags...)
}

func (i *IrysClient) Upload(ctx context.Context, file io.Reader, price *big.Int, tags ...types.Tag) (types.Transaction, error) {
	url := fmt.Sprintf(_uploadPath, i.network, i.currency.GetName())
	b, err := io.ReadAll(file)
	if err != nil {
		return types.Transaction{}, err
	}

	if price == nil {
		price, err = i.GetPrice(ctx, len(b))
		if err != nil {
			return types.Transaction{}, err
		}
	}

	balance, err := i.GetBalance(ctx)
	if err != nil {
		return types.Transaction{}, err
	}

	if balance.Cmp(price) < 0 {
		return types.Transaction{}, errors.ErrBalanceIsLow
	}

	return i.upload(ctx, url, b, tags...)
}

func (i *IrysClient) Download(ctx context.Context, hash string) (*types.File, error) {
	url := fmt.Sprintf(_downloadPath, i.network, hash)

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if err := statusCheck(resp); err != nil {
			return nil, err
		}

		return &types.File{
			Data:          resp.Body,
			Header:        resp.Header,
			ContentLength: resp.ContentLength,
			ContentType:   resp.Header.Get("Content-Type"),
		}, nil
	}
}

func (i *IrysClient) GetMetaData(ctx context.Context, hash string) (types.Transaction, error) {
	url := fmt.Sprintf(_txPath, i.network, hash)

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return types.Transaction{}, err
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return types.Transaction{}, err
	}

	defer resp.Body.Close()

	select {
	case <-ctx.Done():
		return types.Transaction{}, ctx.Err()
	default:
		if err := statusCheck(resp); err != nil {
			return types.Transaction{}, err
		}

		return decodeBody[types.Transaction](resp.Body)
	}
}

func (i *IrysClient) upload(ctx context.Context, url string, payload []byte, tags ...types.Tag) (types.Transaction, error) {
	tags = addContentType(http.DetectContentType(payload), tags...)

	dataItem := types.BundleItem{
		Data: types.Base64String(payload),
		Tags: tags,
	}

	if err := dataItem.Sign(i.currency.GetSinger()); err != nil {
		return types.Transaction{}, err
	}

	reader, err := dataItem.Reader()
	if err != nil {
		return types.Transaction{}, err
	}

	b, err := io.ReadAll(reader)
	if err != nil {
		return types.Transaction{}, err
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return types.Transaction{}, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := i.client.Do(req)
	if err != nil {
		return types.Transaction{}, err
	}
	defer resp.Body.Close()

	select {
	case <-ctx.Done():
		return types.Transaction{}, ctx.Err()
	default:
		if err := statusCheck(resp); err != nil {
			return types.Transaction{}, err
		}

		return decodeBody[types.Transaction](resp.Body)
	}
}
