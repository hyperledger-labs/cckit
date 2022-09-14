package testdata

import (
	"github.com/hyperledger-labs/cckit/extensions/debug"
	"github.com/hyperledger-labs/cckit/extensions/owner"
	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
	"github.com/hyperledger-labs/cckit/router/param/defparam"
	"github.com/hyperledger-labs/cckit/state"
	m "github.com/hyperledger-labs/cckit/state/mapping"
	"github.com/hyperledger-labs/cckit/state/mapping/testdata/schema"
)

func NewSliceIdCC() *router.Chaincode {
	r := router.New(`complexId`)
	debug.AddHandlers(r, `debug`, owner.Only)

	// Mappings for chaincode state
	r.Use(m.MapStates(m.StateMappings{}.
		//key will be <`EntityWithSliceId`, {Id[0]}, {Id[1]},... {Id[len(Id)-1]} >
		Add(&schema.EntityWithSliceId{}, m.PKeyId())))

	r.Init(owner.InvokeSetFromCreator)

	r.Group(`entity`).
		Invoke(`List`, func(c router.Context) (interface{}, error) {
			return c.State().List(&schema.EntityWithSliceId{})
		}).
		Invoke(`Get`, func(c router.Context) (interface{}, error) {
			return c.State().Get(&schema.EntityWithSliceId{Id: state.StringsIdFromStr(c.ParamString(`Id`))})
		}, param.String(`Id`)).
		Invoke(`Insert`, func(c router.Context) (interface{}, error) {
			return nil, c.State().Insert(c.Param())
		}, defparam.Proto(&schema.EntityWithSliceId{}))

	return router.NewChaincode(r)
}
