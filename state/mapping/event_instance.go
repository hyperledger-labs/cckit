package mapping

import (
	"github.com/hyperledger-labs/cckit/serialize"
)

type (
	EventInstance struct {
		instance         interface{}
		eventMapper      EventMapper
		toBytesConverter serialize.ToBytesConverter
	}
)

func NewEventInstance(instance interface{}, eventMapper EventMapper, toBytesConverter serialize.ToBytesConverter) (*EventInstance, error) {
	return &EventInstance{
		instance:         instance,
		eventMapper:      eventMapper,
		toBytesConverter: toBytesConverter,
	}, nil
}

func (ei EventInstance) Name() (string, error) {
	return ei.eventMapper.Name(ei.instance)
}

func (ei EventInstance) ToBytes() ([]byte, error) {
	return ei.toBytesConverter.ToBytesFrom(ei.instance)
}
