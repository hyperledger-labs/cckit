package serialize

import (
	"errors"
	"fmt"
)

type (
	GenericSerializer struct {
	}

	StringSerializer struct {
	}

	BinarySerializer struct {
	}

	JSONSerializer struct {
	}

	//ProtoSerializer struct {
	//}
)

var (
	DefaultSerializer = &GenericSerializer{}
	KeySerializer     = &StringSerializer{}

	ErrOnlyStringSupported = errors.New(`only string supported`)
)

func (g *GenericSerializer) ToBytesFrom(entry interface{}) ([]byte, error) {
	return toBytes(entry)
}

func (g *GenericSerializer) FromBytesTo(serialized []byte, target interface{}) (interface{}, error) {
	return fromBytes(serialized, target)
}

func (g *StringSerializer) ToBytesFrom(entry interface{}) ([]byte, error) {
	switch v := entry.(type) {
	case string:
		return []byte(v), nil
	}
	return nil, ErrOnlyStringSupported
}

func (g *StringSerializer) FromBytesTo(serialized []byte, target interface{}) (interface{}, error) {
	switch t := target.(type) {
	case string:
		return string(serialized), nil
	default:
		return nil, fmt.Errorf(`type=%s: %w`, t, ErrOnlyStringSupported)
	}
}

//func (js *JSONSerializer) ToBytes(entry interface{}) ([]byte, error) {
//	return json.Marshal(entry)
//}
//
//func (js *JSONSerializer) FromBytes(serialized []byte, target interface{}) (interface{}, error) {
//	return JSONUnmarshalPtr(serialized, target)
//}

//func (ps *ProtoSerializer) ToBytes(entry interface{}) ([]byte, error) {
//	return proto.Marshal(entry.(proto.Message))
//}
//
//func (ps *ProtoSerializer) FromBytes(serialized []byte, target interface{}) (interface{}, error) {
//	return convert.FromBytes(serialized, target)
//}
