package types

import (
	"io"
	"math/big"
	"net/http"
)

type NodeInfo struct {
	Version   string            `json:"version"`
	Addresses map[string]string `json:"addresses"`
	Gateway   string            `json:"gateway"`
}

type BalanceResponse struct {
	Balance string `json:"balance"`
}

type TxToBalanceRequest struct {
	TxId string `json:"tx_id"`
}

type TopUpConfirmationResponse struct {
	Confirmed bool `json:"confirmed"`
}

type TopUpConfirmation struct {
	Confirmed bool     `json:"confirmed"`
	Hash      string   `json:"hash"`
	Balance   *big.Int `json:"balance"`
}

type Transaction struct {
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Address   string `json:"address"`
	Owner     string `json:"owner"`
	Signature string `json:"signature"`
	Target    string `json:"target"`
	Tags      []Tag  `json:"tags"`
	Anchor    string `json:"anchor"`
	DataSize  string `json:"data_size"`
	RawSize   string `json:"raw_size"`
}

type File struct {
	Data          io.ReadCloser
	Header        http.Header
	ContentLength int64
	ContentType   string
}

type Chunk struct {
	ID     string
	Offset int64
	Data   []byte
}

type Job struct {
	Chunk Chunk
	Index int
}

type ChunkResponse struct {
	ID  string
	Min int
	Max int
}

type Receipt struct {
	Signature      string `json:"signature"`
	Timestamp      int64  `json:"timestamp"`
	Version        string `json:"version"`
	DeadlineHeight int    `json:"deadlineHeight"`
}

type ReceiptResponse struct {
	Data struct {
		Transactions struct {
			Edges []struct {
				Node struct {
					Receipt struct {
						Signature      string `json:"signature"`
						Timestamp      int64  `json:"timestamp"`
						Version        string `json:"version"`
						DeadlineHeight int    `json:"deadlineHeight"`
					} `json:"receipt"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"transactions"`
	} `json:"data"`
}

type ChunkInfoResponse struct {
	Chunks []int `json:"chunks"`
	Total  int   `json:"total"`
}

func (b BalanceResponse) ToBigInt() *big.Int {
	bInt := new(big.Int)

	n, ok := bInt.SetString(b.Balance, 10)
	if !ok {
		return bInt
	}

	return n
}
