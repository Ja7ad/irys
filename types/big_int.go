package types

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type BigInt struct {
	big.Int
	Valid bool
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	var n json.Number
	err := json.Unmarshal(p, &n)
	if err != nil {
		return err
	}

	var z big.Int
	_, ok := z.SetString(string(n), 10)
	if !ok {
		return fmt.Errorf("not a valid big integer: %s", p)
	}
	b.Int = z
	b.Valid = true
	return nil
}
