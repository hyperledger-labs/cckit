package mapping

import (
	"github.com/hyperledger-labs/cckit/serialize"
	"github.com/hyperledger-labs/cckit/state"
)

type (
	StateInstance struct {
		// instance can be instance itself or key for instance
		// key can be proto or Key ( []string )
		instance         interface{}
		stateMapper      StateMapper
		toBytesConverter serialize.ToBytesConverter
	}
)

func NewStateInstance(instance interface{}, stateMapper StateMapper, toBytesConverter serialize.ToBytesConverter) *StateInstance {
	return &StateInstance{
		instance:         instance,
		stateMapper:      stateMapper,
		toBytesConverter: toBytesConverter,
	}
}

func (si *StateInstance) Key() (state.Key, error) {
	switch instance := si.instance.(type) {
	case []string:
		return instance, nil
	default:
		return si.stateMapper.PrimaryKey(instance)
	}
}

func (si *StateInstance) Keys() ([]state.KeyValue, error) {
	return si.stateMapper.Keys(si.instance, si.toBytesConverter)
}

func (si *StateInstance) ToBytes() ([]byte, error) {
	return si.toBytesConverter.ToBytesFrom(si.instance)
}

func (si *StateInstance) Mapper() StateMapper {
	return si.stateMapper
}
