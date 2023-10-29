package irys

import (
	"context"
	"fmt"
	"github.com/Ja7ad/irys/currency"
	"github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/types"
	"github.com/Ja7ad/irys/utils/logger"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	mu       *sync.Mutex
	client   *retryablehttp.Client
	network  Node
	currency currency.Currency
	contract string
	logging  logger.Logger
	debug    bool
}

type Irys interface {
	// GetPrice return fee base on fileSize in byte for selected currency
	GetPrice(ctx context.Context, fileSize int) (*big.Int, error)

	// BasicUpload file with calculate price and topUp balance base on price (this is slower for upload)
	BasicUpload(ctx context.Context, file io.Reader, tags ...types.Tag) (types.Transaction, error)
	// Upload file with check balance
	Upload(ctx context.Context, file io.Reader, price *big.Int, tags ...types.Tag) (types.Transaction, error)
	// ChunkUpload upload file chunk concurrent
	ChunkUpload(ctx context.Context, file io.Reader, price *big.Int, tags ...types.Tag) (types.Transaction, error)

	// Download get file with header details
	Download(ctx context.Context, hash string) (*types.File, error)
	// GetMetaData get transaction details
	GetMetaData(ctx context.Context, hash string) (types.Transaction, error)

	// GetBalance return current balance in irys node
	GetBalance(ctx context.Context) (*big.Int, error)
	// TopUpBalance top up your balance base on your amount in selected node
	TopUpBalance(ctx context.Context, amount *big.Int) (types.TopUpConfirmation, error)
}

// New create IrysClient object
func New(node Node, currency currency.Currency, debug bool, options ...Option) (Irys, error) {
	irys := new(Client)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	irys.client = retryablehttp.NewClient()
	irys.client.HTTPClient = httpClient

	irys.network = node
	irys.currency = currency
	irys.mu = new(sync.Mutex)

	irys.debug = debug

	irys.client.RetryMax = 5
	irys.client.RetryWaitMin = 1 * time.Second
	irys.client.RetryWaitMax = 30 * time.Second
	irys.client.ErrorHandler = retryablehttp.PassthroughErrorHandler

	for _, opt := range options {
		opt(irys)
	}

	if irys.logging == nil {
		logging, err := logger.New(logger.CONSOLE_HANDLER, logger.Options{
			Development:  false,
			Debug:        false,
			EnableCaller: true,
			SkipCaller:   4,
		})
		if err != nil {
			return nil, err
		}

		irys.logging = logging
	}

	irys.client.Logger = irys.logging

	if !debug {
		irys.client.Logger = nil
	}

	irys.mu.Lock()
	contract, err := irys.getTokenContractAddress(node, currency)
	if err != nil {
		return nil, err
	}
	irys.mu.Unlock()

	irys.contract = contract

	return irys, nil
}

func (i *Client) getTokenContractAddress(node Node, currency currency.Currency) (string, error) {
	r, err := i.client.Get(string(node))
	if err != nil {
		return "", err
	}

	if err := statusCheck(r); err != nil {
		return "", err
	}

	resp, err := decodeBody[types.NodeInfo](r.Body)
	if err != nil {
		return "", err
	}

	if v, ok := resp.Addresses[currency.GetName()]; ok {
		i.debugMsg("set currency address %s base on currency %s", v, currency.GetName())
		return v, nil
	}

	return "", errors.ErrCurrencyIsInvalid
}

func (i *Client) debugMsg(msg string, args ...any) {
	if i.debug {
		i.logging.Debug(fmt.Sprintf(msg, args...))
	}
}
