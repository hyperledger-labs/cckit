package crosscc

import (
	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/state"
	m "github.com/hyperledger-labs/cckit/state/mapping"
)

var (
	StateMappings = m.StateMappings{}.
			Add(&ServiceLocator{},
			m.PKeySchema(&ServiceLocatorId{}),
			m.List(&ServiceLocators{}))

	EventMappings = m.EventMappings{}.
			Add(&ServiceLocatorSet{})
)

func State(ctx router.Context) m.MappedState {
	return m.WrapState(ctx.State(), StateMappings)
}

func Event(ctx router.Context) state.Event {
	return m.WrapEvent(ctx.Event(), EventMappings)
}
