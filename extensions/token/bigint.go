package token

import (
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/golang/protobuf/jsonpb"
)

var ErrConvertStringInt = errors.New("failed to convert string to big int")

func (x *BigInt) BigInt() (*big.Int, error) {
	bigInt, ok := new(big.Int).SetString(x.Data, 10)
	if !ok {
		return nil, ErrConvertStringInt
	}
	return bigInt, nil
}

// MarshalJSON collapse field json representation to string from Data field
func (x *BigInt) MarshalJSON() ([]byte, error) {
	return []byte(strings.Join([]string{`"`, x.GetData(), `"`}, ``)), nil
}

func (x *BigInt) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
	return x.MarshalJSON()
}

type bigIntUnmarshal struct {
	Data string
}

// UnmarshalJSON supports short ( {"amount: "12345"}) and full {"amount: { "data" : "12345" }} BigInt json representation
func (x *BigInt) UnmarshalJSON(data []byte) error {
	str := string(data)
	strLen := len(str)

	if str[0:1] == `{` && str[strLen-1:strLen] == `}` {
		obj := new(bigIntUnmarshal)
		if err := json.Unmarshal(data, obj); err != nil {
			return err
		}
		x.Data = obj.Data
		return nil
	}
	x.Data = str[1 : strLen-1] // remove quotes
	return nil
}

func (x *BigInt) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, data []byte) error {
	return x.UnmarshalJSON(data)
}

func NewBigInt(val *big.Int) *BigInt {
	return &BigInt{
		Data: val.String(),
	}
}

func BigIntSum(a, b *big.Int) *big.Int {
	return new(big.Int).Add(a, b)
}

func NewBigIntSum(a, b *big.Int) *BigInt {
	return NewBigInt(BigIntSum(a, b))
}

func BigIntSub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

func NewBigIntSub(a, b *big.Int) *BigInt {
	return NewBigInt(BigIntSub(a, b))
}
