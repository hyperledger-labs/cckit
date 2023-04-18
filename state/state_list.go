package state

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/hyperledger-labs/cckit/serialize"
	"github.com/hyperledger-labs/cckit/state/schema"
)

type (
	List struct {
		itemTarget interface{}   // type of item
		listTarget interface{}   // type of item list
		list       []interface{} // actual lost
	}
)

// ListItemsField name of field with items in listTarget proto structure setting
const ListItemsField = `Items`

func NewList(config ...interface{}) (*List, error) {
	var (
		itemTarget, listTarget interface{}
	)
	if len(config) > 0 {
		itemTarget = config[0]
	}
	if len(config) > 1 {
		listTarget = config[1]
	}

	return &List{itemTarget: itemTarget, listTarget: listTarget}, nil
}

// Fill state list from iterator
func (sl *List) Fill(
	iter shim.StateQueryIteratorInterface, fromBytesConverter serialize.FromBytesConverter) (list interface{}, err error) {
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return nil, err
		}
		item, err := fromBytesConverter.FromBytesTo(kv.Value, sl.itemTarget)
		if err != nil {
			return nil, fmt.Errorf(`transform list entry: %w`, err)
		}
		sl.list = append(sl.list, item)
	}
	return sl.Get()
}

func (sl *List) Get() (list interface{}, err error) {
	// list type is  proto.Message, with predefined Items field
	if _, isListProto := sl.listTarget.(proto.Message); isListProto {
		customList := proto.Clone(sl.listTarget.(proto.Message)) // create copy of list type proto message
		items := reflect.ValueOf(customList).Elem().FieldByName(`Items`)

		for _, v := range sl.list {
			items.Set(reflect.Append(items, reflect.ValueOf(v)))
		}
		return customList, nil

		// default list proto.Message ( with repeated Any)
	} else if _, isItemProto := sl.itemTarget.(proto.Message); isItemProto {
		defList := &schema.List{}

		for _, item := range sl.list {
			any, err := ptypes.MarshalAny(item.(proto.Message))
			if err != nil {
				return nil, err
			}
			defList.Items = append(defList.Items, any)
		}
		return defList, nil
	}

	return sl.list, nil
}

func (sl *List) AddElementToList(elem interface{}) {
	sl.list = append(sl.list, elem)
}
