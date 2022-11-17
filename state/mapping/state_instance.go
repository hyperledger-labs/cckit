package mapping

import (
	"github.com/hyperledger-labs/cckit/serialize"
	"github.com/hyperledger-labs/cckit/state"
)

type (
	StateInstance struct {
		// StateInstance  can be instance itself or key for instance
		// key can be proto or Key ( []string )
		instance interface{}
		mapper   StateMapper
	}
)

func NewStateInstance(instance interface{}, stateMapper StateMapper) *StateInstance {
	return &StateInstance{
		instance: instance,
		mapper:   stateMapper,
	}
}

func (si *StateInstance) Key() (state.Key, error) {
	switch instance := si.instance.(type) {
	case []string:
		return instance, nil
	default:
		return si.mapper.PrimaryKey(instance)
	}
}

func (si *StateInstance) ToBytes(toBytesConverter serialize.ToBytesConverter) ([]byte, error) {
	return toBytesConverter.ToBytesFrom(si.instance)
}

func (si *StateInstance) Keys() ([]state.KeyValue, error) {
	return si.mapper.Keys(si.instance)
}

func (si *StateInstance) Mapper() StateMapper {
	return si.mapper
}
