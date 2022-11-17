package state

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/hyperledger-labs/cckit/serialize"
)

type (
	// Event interface for working with events in chaincode
	Event interface {
		Set(entry interface{}, value ...interface{}) error
		UseToBytesConverter(serialize.ToBytesConverter) Event
		// 	ToBytesConverter todo: check neediness
		ToBytesConverter() serialize.ToBytesConverter
		UseNameTransformer(StringTransformer) Event
	}

	Namer interface {
		Name() (string, error)
	}
	// NameValue interface combines Name() as ToByter methods - event representation
	NameValue interface {
		Namer
		serialize.ToBytesConverter
	}

	EventImpl struct {
		stub             shim.ChaincodeStubInterface
		nameTransformer  StringTransformer
		toBytesConverter serialize.ToBytesConverter
	}
)

// NewEvent creates wrapper on shim.ChaincodeStubInterface for working with events
func NewEvent(stub shim.ChaincodeStubInterface) *EventImpl {
	return &EventImpl{
		stub:             stub,
		nameTransformer:  NameAsIs,
		toBytesConverter: serialize.DefaultSerializer,
	}
}

func (e *EventImpl) UseToBytesConverter(toBytesConverter serialize.ToBytesConverter) Event {
	e.toBytesConverter = toBytesConverter
	return e
}

func (e *EventImpl) ToBytesConverter() serialize.ToBytesConverter {
	return e.toBytesConverter
}

func (e *EventImpl) UseNameTransformer(nt StringTransformer) Event {
	e.nameTransformer = nt
	return e
}

func (e *EventImpl) Set(entry interface{}, values ...interface{}) error {
	name, value, err := e.ArgNameValue(entry, values)
	if err != nil {
		return err
	}

	nameStr, err := e.nameTransformer(name)
	if err != nil {
		return err
	}

	bb, err := e.toBytesConverter.ToBytesFrom(value)
	if err != nil {
		return err
	}

	return e.stub.SetEvent(nameStr, bb)
}

func (e *EventImpl) ArgNameValue(arg interface{}, values []interface{}) (name string, value interface{}, err error) {
	// name must be
	name, err = NormalizeEventName(arg)
	if err != nil {
		return
	}

	switch len(values) {
	// arg is name and  value
	case 0:
		return name, arg, nil
	case 1:
		return name, values[0], nil
	default:
		return ``, nil, ErrAllowOnlyOneValue
	}
}

func NormalizeEventName(name interface{}) (string, error) {
	switch n := name.(type) {
	case Namer:
		return n.Name()
	case string:
		return n, nil
	}

	return ``, ErrUnableToCreateEventName
}
