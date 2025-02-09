package accounts_test

import (
	intMaxAcc "intmax2-withdrawal/internal/accounts"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeG1CurvePoint(t *testing.T) {
	t.Parallel()

	randomScalar, err := new(fr.Element).SetRandom()
	assert.NoError(t, err)

	point := new(bn254.G1Affine).ScalarMultiplicationBase(randomScalar.BigInt(new(big.Int)))
	encoded := intMaxAcc.EncodeG1CurvePoint(point)
	decoded, err := intMaxAcc.DecodeG1CurvePoint(encoded)
	assert.NoError(t, err)
	assert.Equal(t, point, decoded, "Decoded point is different from the original one")
}

func TestEncodeDecodeG2CurvePoint(t *testing.T) {
	t.Parallel()

	point, err := intMaxAcc.DecodeG2CurvePoint("27f811fe50964adcb0345ddf85dd0e2e913229991b1d2a551df2908e8ccd3bfc2ba7d3c0ce4096f524d22afeba96b6ce95a6357b5336f9cc57dc0cc78fa605e604781cec49a668fc7ec5dc22fd5f9e49e2b594b1ff9b8067c97d2b60d6be6cd0048da9489637392dc5c427d7b5e9b0976158a3f06b58820c90245ad68675b8b4")
	require.NoError(t, err)
	encoded := intMaxAcc.EncodeG2CurvePoint(point)
	decoded, err := intMaxAcc.DecodeG2CurvePoint(encoded)
	assert.NoError(t, err)
	assert.Equal(t, point, decoded, "Decoded point is different from the original one")
}
