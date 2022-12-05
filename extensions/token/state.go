package token

import (
	"errors"

	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/state"
	m "github.com/hyperledger-labs/cckit/state/mapping"
)

var (
	ErrAmountInsuficcient = errors.New(`amount insufficient`)

	StateMappings = m.StateMappings{}.
		//  Create mapping for Balance entity (account model)
		// key will be `Balance`,`{Address}`,`{Path[0]}`..., `{Path[n]`
		Add(&Balance{},
			m.PKeySchema(&BalanceId{}),
			m.List(&Balances{}), // Structure of result for List method
		).
		//  Create mapping for Balance entity
		// key will be `UTXO`,`{Address}`,`{symbol}`, `{group | join (",")}`, `{TXId}`
		Add(&UTXO{},
			//m.PKeyer(func(instance interface{}) (state.Key, error) {
			//	return UTXOKey(instance.(*UTXO)), nil
			//}),
			m.PKeySchema(&UTXOId{}),
			m.List(&UTXOs{}), // Structure of result for List method
		).
		Add(&TokenType{},
			m.PKeySchema(&TokenTypeId{}),
			m.List(&TokenTypes{})).
		Add(&TokenGroup{},
			m.PKeySchema(&TokenGroupId{}),
			m.List(&TokenGroups{})).
		Add(&Config{},
			m.WithConstPKey())

	EventMappings = m.EventMappings{}.
			Add(&Transferred{}).
			Add(&TransferredBatch{}).
			Add(&TokenTypeCreated{}).
			Add(&TokenGroupCreated{})
)

func UTXOKey(utxo *UTXO) state.Key {
	return UTXOKeyBase(utxo).Append(state.Key{utxo.TxId})
}

//func UTXOKeyBase(utxo *UTXO) state.Key {
//	return state.Key{utxo.Address, utxo.Symbol, strings.Join(utxo.Group, `,`)}
//}

func UTXOKeyBase(utxo *UTXO) state.Key {
	return state.Key{utxo.Symbol, utxo.Group, utxo.Address}
}

// State with chaincode mappings
func State(ctx router.Context) m.MappedState {
	return m.WrapState(ctx.State(), StateMappings)
}

// Event with chaincode mappings
func Event(ctx router.Context) state.Event {
	return m.WrapEvent(ctx.Event(), EventMappings)
}
