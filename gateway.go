package irys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/types"
	"github.com/ethereum/go-ethereum/crypto"
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

func (i *Irys) GetPrice(fileSize int) (*big.Int, error) {
	url := fmt.Sprintf(_pricePath, i.network, i.currency.GetName(), fileSize)
	resp, err := i.client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(b))
	}

	return decodeBody[*big.Int](resp.Body)
}

func (i *Irys) GetBalance() (*big.Int, error) {
	pbKey := i.currency.GetPublicKey()
	url := fmt.Sprintf(_getBalance, i.network, crypto.PubkeyToAddress(*pbKey).Hex())

	resp, err := i.client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(b))
	}

	b, err := decodeBody[types.BalanceResponse](resp.Body)
	if err != nil {
		return nil, err
	}

	return b.ToBigInt(), nil
}

func (i *Irys) TopUpBalance(ctx context.Context, amount *big.Int) (types.TopUpConfirmation, error) {
	urlConfirm := fmt.Sprintf(_sendTxToBalance, i.network)

	hash, err := i.createTx(ctx, amount)
	if err != nil {
		return types.TopUpConfirmation{}, err
	}

	for {
		select {
		case <-ctx.Done():
			return types.TopUpConfirmation{}, ctx.Err()
		default:
			b, err := json.Marshal(&types.TxToBalanceRequest{
				TxId: hash,
			})
			if err != nil {
				return types.TopUpConfirmation{}, err
			}
			resp, err := i.client.Post(urlConfirm, "application/json", bytes.NewBuffer(b))
			if err != nil {
				return types.TopUpConfirmation{}, err
			}

			if resp.StatusCode != http.StatusOK {
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					return types.TopUpConfirmation{}, err
				}
				return types.TopUpConfirmation{}, fmt.Errorf("%d: %s", resp.StatusCode, string(b))
			}

			confirm, err := decodeBody[types.TopUpConfirmationResponse](resp.Body)
			if err != nil {
				return types.TopUpConfirmation{}, err
			}

			if confirm.Confirmed {
				balance, err := i.GetBalance()
				if err != nil {
					return types.TopUpConfirmation{}, err
				}
				return types.TopUpConfirmation{
					Confirmed: true,
					Hash:      hash,
					Balance:   balance,
				}, nil
			}
		}
	}
}

func (i *Irys) BasicUpload(ctx context.Context, file io.Reader, tags ...types.Tag) (types.Transaction, error) {
	url := fmt.Sprintf(_uploadPath, i.network, i.currency.GetName())
	b, err := io.ReadAll(file)
	if err != nil {
		return types.Transaction{}, err
	}

	price, err := i.GetPrice(len(b))
	if err != nil {
		return types.Transaction{}, err
	}

	balance, err := i.GetBalance()
	if err != nil {
		return types.Transaction{}, err
	}

	if balance.Cmp(price) < 0 {
		_, err := i.TopUpBalance(ctx, price)
		if err != nil {
			return types.Transaction{}, err
		}
	}

	return i.upload(url, b, tags...)
}

func (i *Irys) Upload(ctx context.Context, file io.Reader, price *big.Int, tags ...types.Tag) (types.Transaction, error) {
	url := fmt.Sprintf(_uploadPath, i.network, i.currency.GetName())
	b, err := io.ReadAll(file)
	if err != nil {
		return types.Transaction{}, err
	}

	if price == nil {
		price, err = i.GetPrice(len(b))
		if err != nil {
			return types.Transaction{}, err
		}
	}

	balance, err := i.GetBalance()
	if err != nil {
		return types.Transaction{}, err
	}

	if balance.Cmp(price) < 0 {
		return types.Transaction{}, errors.ErrBalanceIsLow
	}

	return i.upload(url, b, tags...)
}

func (i *Irys) Download(hash string) (*types.File, error) {
	url := fmt.Sprintf(_downloadPath, i.network, hash)
	resp, err := i.client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(b))
	}

	return &types.File{
		Data:          resp.Body,
		Header:        resp.Header,
		ContentLength: resp.ContentLength,
		ContentType:   resp.Header.Get("Content-Type"),
	}, nil
}

func (i *Irys) GetMetaData(hash string) (types.Transaction, error) {
	url := fmt.Sprintf(_txPath, i.network, hash)
	resp, err := i.client.Get(url)
	if err != nil {
		return types.Transaction{}, err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return types.Transaction{}, err
		}
		return types.Transaction{}, fmt.Errorf("%d: %s", resp.StatusCode, string(b))
	}

	return decodeBody[types.Transaction](resp.Body)
}

func (i *Irys) upload(url string, payload []byte, tags ...types.Tag) (types.Transaction, error) {
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

	resp, err := i.client.Post(url, "application/octet-stream", bytes.NewBuffer(b))
	if err != nil {
		return types.Transaction{}, err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return types.Transaction{}, err
		}
		return types.Transaction{}, fmt.Errorf("%d: %s", resp.StatusCode, string(b))
	}

	return decodeBody[types.Transaction](resp.Body)
}
