package serialize

import (
	"errors"
)

var (
	//ErrBytesToSerializeEmpty = errors.New(`bytes to serialize empty`)

	// ErrUnableToConvertNilToStruct - nil cannot be converted to struct
	ErrUnableToConvertNilToStruct = errors.New(`unable to convert nil to [struct,array,slice,ptr]`)
	// ErrUnableToConvertValueToStruct - value  cannot be converted to struct
	ErrUnableToConvertValueToStruct = errors.New(`unable to convert value to struct`)
)
