package token

import (
	"errors"
	"math/big"
)

var ErrConvertStringInt = errors.New("failed to convert string to big int")

func IntVal(x string) (*big.Int, error) {
	X, ok := new(big.Int).SetString(x, 10)
	if !ok {
		return nil, ErrConvertStringInt
	}
	return X, nil
}
