package irys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	errs "github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/types"
	retry "github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"math/big"
	"net/http"
	"strings"
)

const (
	_pricePath       = "%s/price/%s/%v"
	_uploadPath      = "%s/tx/%s"
	_txPath          = "%s/tx/%s"
	_downloadPath    = "%s/%s"
	_sendTxToBalance = "%s/account/balance/matic"
	_getBalance      = "%s/account/balance/matic?address=%s"
	_chunkUpload     = "%s/chunks/%s/%v/%v"
	_graphql         = "%s/graphql"
)

func (c *Client) GetPrice(ctx context.Context, fileSize int) (*big.Int, error) {
	url := fmt.Sprintf(_pricePath, c.network, c.currency.GetName(), fileSize)
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
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

func (c *Client) GetBalance(ctx context.Context) (*big.Int, error) {
	pbKey := c.currency.GetPublicKey()
	url := fmt.Sprintf(_getBalance, c.network, crypto.PubkeyToAddress(*pbKey).Hex())

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
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

func (c *Client) TopUpBalance(ctx context.Context, amount *big.Int) (types.TopUpConfirmation, error) {
	urlConfirm := fmt.Sprintf(_sendTxToBalance, c.network)

	hash, err := c.createTx(ctx, amount)
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

	resp, err := c.client.Do(req)
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
				balance, err = c.GetBalance(ctx)
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

func (c *Client) Download(ctx context.Context, txId string) (*types.File, error) {
	url := fmt.Sprintf(_downloadPath, _defaultGateway, txId)

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
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

func (c *Client) GetMetaData(ctx context.Context, txId string) (types.Transaction, error) {
	url := fmt.Sprintf(_txPath, _defaultGateway, txId)

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return types.Transaction{}, err
	}

	resp, err := c.client.Do(req)
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

func (c *Client) BasicUpload(ctx context.Context, file io.Reader, tags ...types.Tag) (types.Transaction, error) {
	url := fmt.Sprintf(_uploadPath, c.network, c.currency.GetName())
	b, err := io.ReadAll(file)
	if err != nil {
		return types.Transaction{}, err
	}

	price, err := c.GetPrice(ctx, len(b))
	if err != nil {
		return types.Transaction{}, err
	}
	c.debugMsg("[BasicUpload] get price %s", price.String())

	balance, err := c.GetBalance(ctx)
	if err != nil {
		return types.Transaction{}, err
	}
	c.debugMsg("[BasicUpload] get balance %s", balance.String())

	if balance.Cmp(price) < 0 {
		_, err := c.TopUpBalance(ctx, price)
		if err != nil {
			return types.Transaction{}, err
		}
		c.debugMsg("[BasicUpload] topUp balance")
	}

	return c.upload(ctx, url, b, tags...)
}

func (c *Client) GetReceipt(ctx context.Context, txId string) (types.Receipt, error) {
	url := fmt.Sprintf(_graphql, c.network)

	body := strings.NewReader(fmt.Sprintf("{\"query\":\"query {\\n      "+
		"transactions(ids: [\\\"%s\\\"]) {\\n        edges {\\n         "+
		" node {\\n            receipt {\\n              signature\\n              "+
		"timestamp\\n              version\\n              deadlineHeight\\n           "+
		" }\\n          }\\n        }\\n      }\\n    }\",\"variables\":{}}", txId))

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return types.Receipt{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return types.Receipt{}, err
	}

	defer resp.Body.Close()

	select {
	case <-ctx.Done():
		return types.Receipt{}, ctx.Err()
	default:
		if err := statusCheck(resp); err != nil {
			return types.Receipt{}, err
		}

		response, err := decodeBody[types.ReceiptResponse](resp.Body)
		if err != nil {
			return types.Receipt{}, err
		}

		if response.Data.Transactions.Edges != nil {
			return response.Data.Transactions.Edges[0].Node.Receipt, nil
		}

		return types.Receipt{}, nil
	}
}

func (c *Client) Upload(ctx context.Context, file io.Reader, price *big.Int, tags ...types.Tag) (types.Transaction, error) {
	url := fmt.Sprintf(_uploadPath, c.network, c.currency.GetName())
	b, err := io.ReadAll(file)
	if err != nil {
		return types.Transaction{}, err
	}

	if price == nil {
		price, err = c.GetPrice(ctx, len(b))
		if err != nil {
			return types.Transaction{}, err
		}
		c.debugMsg("[Upload] get price %s", price.String())
	}

	balance, err := c.GetBalance(ctx)
	if err != nil {
		return types.Transaction{}, err
	}
	c.debugMsg("[Upload] get balance %s", balance.String())

	if balance.Cmp(price) < 0 {
		return types.Transaction{}, errs.ErrBalanceIsLow
	}

	return c.upload(ctx, url, b, tags...)
}

func (c *Client) upload(ctx context.Context, url string, file []byte, tags ...types.Tag) (types.Transaction, error) {
	b, err := signFile(file, c.currency.GetSinger(), true, tags...)
	if err != nil {
		return types.Transaction{}, err
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return types.Transaction{}, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	c.debugMsg("[Upload] create upload request")

	resp, err := c.client.Do(req)
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
