package serialize

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/encoding/protojson"
)

type (
	ProtoSerializer interface {
		ToBytes(proto.Message) ([]byte, error)
		FromBytes(serialized []byte, target proto.Message) (proto.Message, error)
	}

	BinaryProtoSerializer struct {
	}

	JSONProtoSerializer struct {
		UseProtoNames bool // to use proto field name instead of lowerCamelCase name in JSON field names
	}
)

func (ps *BinaryProtoSerializer) ToBytes(entry proto.Message) ([]byte, error) {
	return BinaryProtoMarshal(entry)
}

func (ps *BinaryProtoSerializer) FromBytes(serialized []byte, target proto.Message) (proto.Message, error) {
	return BinaryProtoUnmarshal(serialized, target)
}

func (js *JSONProtoSerializer) ToBytes(entry proto.Message) ([]byte, error) {
	mo := &protojson.MarshalOptions{UseProtoNames: js.UseProtoNames}
	return JSONProtoMarshal(entry, mo)
}

func (js *JSONProtoSerializer) FromBytes(serialized []byte, target proto.Message) (proto.Message, error) {
	return JSONProtoUnmarshal(serialized, target)
}

func BinaryProtoMarshal(entry proto.Message) ([]byte, error) {
	return proto.Marshal(proto.Clone(entry))
}

func JSONProtoMarshal(entry proto.Message, mo *protojson.MarshalOptions) ([]byte, error) {
	// return protojson.Marshal(proto.Clone(entry))
	// mo := protojson.MarshalOptions{UseProtoNames: true}
	return mo.Marshal(proto.Clone(entry))
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
	err = protojson.Unmarshal(json, msg)

	if err != nil {
		return nil, fmt.Errorf(`json proto unmarshal: %w`, err)
	}
	return msg, nil
}
