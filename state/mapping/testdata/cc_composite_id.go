package testdata

import (
	"github.com/hyperledger-labs/cckit/extensions/debug"
	"github.com/hyperledger-labs/cckit/extensions/owner"
	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param/defparam"
	"github.com/hyperledger-labs/cckit/serialize"
	"github.com/hyperledger-labs/cckit/state"
	"github.com/hyperledger-labs/cckit/state/mapping"
	"github.com/hyperledger-labs/cckit/state/mapping/testdata/schema"
)

var (
	EntityCompositeIdNamespace        = state.Key{`entity-composite-id`}
	EntityWithCompositeIdStateMapping = mapping.StateMappings{}.Add(&schema.EntityWithCompositeId{},
		// explicit set namespace for primary key, otherwise namespace will be based of schema type string representation
		mapping.WithNamespace(EntityCompositeIdNamespace),
		//  schema for Primary Key
		mapping.PKeySchema(&schema.EntityCompositeId{}),
		mapping.List(&schema.EntityWithCompositeIdList{}))
)

func NewCompositeIdCC(serializer serialize.Serializer) *router.Chaincode {
	r := router.New("composite_id", router.WithSerializer(serializer))
	r.Use(mapping.MapStates(EntityWithCompositeIdStateMapping))

	r.Use(mapping.MapEvents(mapping.EventMappings{}.
		Add(&schema.CreateEntityWithCompositeId{}).
		Add(&schema.UpdateEntityWithCompositeId{})))

	r.Init(owner.InvokeSetFromCreator)
	debug.AddHandlers(r, "debug", owner.Only)

	r.
		Query(ListFunc, queryListComposite).
		Query(GetFunc, queryByIdComposite, defparam.Proto(&schema.EntityCompositeId{})).
		Invoke(CreateFunc, invokeCreateComposite, defparam.Proto(&schema.CreateEntityWithCompositeId{})).
		Invoke(UpdateFunc, invokeUpdateComposite, defparam.Proto(&schema.UpdateEntityWithCompositeId{})).
		Invoke(DeleteFunc, invokeDeleteComposite, defparam.Proto(&schema.EntityCompositeId{}))

	return router.NewChaincode(r)
}

func queryByIdComposite(c router.Context) (interface{}, error) {
	return c.State().Get(c.Param().(*schema.EntityCompositeId))
}

func queryListComposite(c router.Context) (interface{}, error) {
	return c.State().List(&schema.EntityWithCompositeId{})
}

func invokeCreateComposite(c router.Context) (interface{}, error) {
	create := c.Param().(*schema.CreateEntityWithCompositeId)
	entity := &schema.EntityWithCompositeId{
		IdFirstPart:  create.IdFirstPart,
		IdSecondPart: create.IdSecondPart,
		IdThirdPart:  create.IdThirdPart,
		Name:         create.Name,
		Value:        create.Value,
	}

	if err := c.Event().Set(create); err != nil {
		return nil, err
	}

	return entity, c.State().Insert(entity)
}

func invokeUpdateComposite(c router.Context) (interface{}, error) {
	update := c.Param().(*schema.UpdateEntityWithCompositeId)
	entity, _ := c.State().Get(
		&schema.EntityCompositeId{
			IdFirstPart:  update.IdFirstPart,
			IdSecondPart: update.IdSecondPart,
			IdThirdPart:  update.IdThirdPart,
		},
		&schema.EntityWithCompositeId{})

	e := entity.(*schema.EntityWithCompositeId)

	e.Name = update.Name
	e.Value = update.Value

	if err := c.Event().Set(update); err != nil {
		return nil, err
	}

	return e, c.State().Put(e)
}

func invokeDeleteComposite(c router.Context) (interface{}, error) {
	return nil, c.State().Delete(c.Param().(*schema.EntityCompositeId))
}
