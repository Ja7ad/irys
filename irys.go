package irys

import (
	"context"
	"github.com/Ja7ad/irys/currency"
	"github.com/Ja7ad/irys/errors"
	"github.com/Ja7ad/irys/types"
	"io"
	"math/big"
	"net/http"
	"sync"
)

type Irys struct {
	mu       *sync.Mutex
	client   *http.Client
	network  Node
	currency currency.Currency
	contract string
}

type Gateway interface {
	// GetPrice return fee base on fileSize in byte for selected currency
	GetPrice(fileSize int) (*big.Int, error)

	// BasicUpload file with calculate price and topUp balance base on price (this is slower for upload)
	BasicUpload(ctx context.Context, file io.Reader, tags ...types.Tag) (types.Transaction, error)
	// Upload file with check balance
	Upload(ctx context.Context, file io.Reader, price *big.Int, tags ...types.Tag) (types.Transaction, error)

	// Download get file with header details
	Download(hash string) (*types.File, error)
	// GetMetaData get transaction details
	GetMetaData(hash string) (types.Transaction, error)

	// GetBalance return current balance in irys node
	GetBalance() (*big.Int, error)
	// TopUpBalance top up your balance base on your amount in selected node
	TopUpBalance(ctx context.Context, amount *big.Int) (types.TopUpConfirmation, error)
}

// New create Irys object
func New(node Node, currency currency.Currency, options ...Option) (Gateway, error) {
	irys := new(Irys)
	irys.client = http.DefaultClient
	irys.network = node
	irys.currency = currency
	irys.mu = new(sync.Mutex)

	for _, opt := range options {
		opt(irys)
	}

	irys.mu.Lock()
	defer irys.mu.Unlock()
	contract, err := irys.getTokenContractAddress(node, currency)
	if err != nil {
		return nil, err
	}

	irys.contract = contract

	return irys, nil
}

func (i *Irys) getTokenContractAddress(node Node, currency currency.Currency) (string, error) {
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
		return v, nil
	}

	return "", errors.ErrCurrencyIsInvalid
}
