package token

import (
	"errors"
	"math/big"
)

var ErrConvertStringInt = errors.New("failed to convert string to big int")

func (x *Decimal) BigInt() (*big.Int, error) {
	bigInt, ok := new(big.Int).SetString(x.Value, 10)
	if !ok {
		return nil, ErrConvertStringInt
	}
	return bigInt, nil
}

func NewDecimal(val *big.Int, scale ...int32) *Decimal {
	d := &Decimal{
		Value: val.String(),
	}

	if len(scale) > 0 {
		d.Scale = scale[0]
	}

	return d
}

func BigIntSum(a, b *big.Int) *big.Int {
	return new(big.Int).Add(a, b)
}

func BigIntSumAsDecimal(a, b *big.Int, scale ...int32) *Decimal {
	return NewDecimal(BigIntSum(a, b), scale...)
}

func BigIntSub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

func BigIntSubAsDecimal(a, b *big.Int, scale ...int32) *Decimal {
	return NewDecimal(BigIntSub(a, b), scale...)
}
