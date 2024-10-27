package types

import (
	"encoding/binary"
	"intmax2-withdrawal/internal/hash/goldenposeidon"
	"math/big"
)

const (
	NumPublicKeyBytes   = 32
	PublicKeySenderType = "PUBLIC_KEY"

	NumAccountIDBytes   = 5
	AccountIDSenderType = "ACCOUNT_ID"

	NumOfSenders                  = 128
	numFlagBytes                  = 16
	numBaseFieldOrderBytes        = 32
	numG1PointLimbs               = 2
	numG2PointLimbs               = 4
	defaultAccountID       uint64 = 0
	dummyAccountID         uint64 = 1
	int8Key                       = 8
	int32Key                      = 32
	int10Key                      = 10
	int64Key                      = 64
	int128Key                     = 128
)

type PoseidonHashOut = goldenposeidon.PoseidonHashOut

func BigIntToBytes32BeArray(bi *big.Int) [int32Key]byte {
	biBytes := bi.Bytes()
	var result [int32Key]byte
	copy(result[int32Key-len(biBytes):], biBytes)
	return result
}

type Bytes32 [int8Key]uint32

func (b *Bytes32) FromBytes(bytes []byte) {
	for i := 0; i < int8Key; i++ {
		b[i] = binary.BigEndian.Uint32(bytes[i*int4Key : (i+1)*int4Key])
	}
}

func (b *Bytes32) Bytes() []byte {
	bytes := make([]byte, int8Key*int4Key)
	for i := 0; i < int8Key; i++ {
		binary.BigEndian.PutUint32(bytes[i*int4Key:(i+1)*int4Key], b[i])
	}

	return bytes
}

func Uint32SliceToBytes(v []uint32) []byte {
	const int4Key = 4

	buf := make([]byte, len(v)*int4Key)
	for i, n := range v {
		binary.BigEndian.PutUint32(buf[i*int4Key:], n)
	}

	return buf
}
