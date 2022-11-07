package serialize

import (
	"errors"
	"fmt"
)

var (
	ErrToBytesConverterIsNil = errors.New(`to bytes converter is nil`)
)

// ArgsToBytes converts func arguments to bytes
func ArgsToBytes(iargs []interface{}, toBytesConverter ToBytesConverter) (aa [][]byte, err error) {
	if toBytesConverter == nil {
		return nil, ErrToBytesConverterIsNil
	}
	args := make([][]byte, len(iargs))

	for i, arg := range iargs {
		val, err := toBytesConverter.ToBytesFrom(arg)
		if err != nil {
			return nil, fmt.Errorf(`convert arg[%d]: %w`, i, err)
		}
		args[i] = val
	}

	return args, nil
}
