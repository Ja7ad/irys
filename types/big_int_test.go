package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBigIntMarshalUnmarshal(t *testing.T) {
	n := BigInt{*big.NewInt(123), true}

	json, err := n.MarshalJSON()
	require.NoError(t, err)

	m := BigInt{}
	err = m.UnmarshalJSON(json)
	require.NoError(t, err)

	require.Equal(t, n, m)
}
