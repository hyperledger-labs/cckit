package serialize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type (
	ProtoSerializer interface {
		ToBytes(proto.Message) ([]byte, error)
		FromBytes(serialized []byte, target proto.Message) (proto.Message, error)
	}

	BinaryProtoSerializer struct {
	}

	JSONProtoSerializer struct {
	}
)

func (ps *BinaryProtoSerializer) ToBytes(entry proto.Message) ([]byte, error) {
	return BinaryProtoMarshal(entry)
}

func (ps *BinaryProtoSerializer) FromBytes(serialized []byte, target proto.Message) (proto.Message, error) {
	return BinaryProtoUnmarshal(serialized, target)
}

func (js *JSONProtoSerializer) ToBytes(entry proto.Message) ([]byte, error) {
	return JSONProtoMarshal(entry)
}

func (js *JSONProtoSerializer) FromBytes(serialized []byte, target proto.Message) (proto.Message, error) {
	return JSONProtoUnmarshal(serialized, target)
}

func BinaryProtoMarshal(entry proto.Message) ([]byte, error) {
	return proto.Marshal(proto.Clone(entry))
}

func JSONProtoMarshal(entry proto.Message) ([]byte, error) {
	return json.Marshal(proto.Clone(entry))
}

// BinaryProtoUnmarshal r unmarshalls []byte as proto.Message to pointer, and returns value pointed to
func BinaryProtoUnmarshal(bb []byte, messageType proto.Message) (message proto.Message, err error) {
	msg := proto.Clone(messageType)
	err = proto.Unmarshal(bb, msg)
	if err != nil {
		return nil, fmt.Errorf(`unmarshal to proto=%s: %w`, reflect.TypeOf(messageType), err)
	}
	return msg, nil
}

func JSONProtoUnmarshal(json []byte, messageType proto.Message) (message proto.Message, err error) {
	msg := proto.Clone(messageType)
	err = jsonpb.Unmarshal(bytes.NewReader(json), msg)
	if err != nil {
		return nil, fmt.Errorf(`json proto unmarshal: %w`, err)
	}
	return msg, nil
}
