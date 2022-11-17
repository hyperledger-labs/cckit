package mapping

import (
	"github.com/hyperledger-labs/cckit/serialize"
)

type (
	EventInstance struct {
		instance interface{}
		mapper   EventMapper
	}
)

func NewEventInstance(instance interface{}, mapper EventMapper) (*EventInstance, error) {
	return &EventInstance{
		instance: instance,
		mapper:   mapper,
	}, nil
}

func (ei *EventInstance) Name() (string, error) {
	return ei.mapper.Name(ei.instance)
}

func (ei *EventInstance) ToBytes(toBytesConverter serialize.ToBytesConverter) ([]byte, error) {
	return toBytesConverter.ToBytesFrom(ei.instance)
}
