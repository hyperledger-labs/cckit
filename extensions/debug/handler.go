package debug

import (
	"github.com/pkg/errors"

	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
	"github.com/hyperledger-labs/cckit/state"
)

const (
	InvokeStateCleanFunc  = `StateClean`
	QueryStateKeysFunc    = `StateKeys`
	QueryStateGetFunc     = `StateGet`
	InvokeStatePutFunc    = `StatePut`
	InvokeStateDeleteFunc = `StateDelete`
)

var (
	// KeyParam parameter for get, put, delete data from state
	KeyParam = param.Strings(`key`)

	// PrefixParam parameter
	PrefixParam = param.String(`prefix`)

	// PrefixesParams parameter
	PrefixesParam = param.Strings(`prefixes`)

	// ValueParam  parameter for putting value in state
	ValueParam = param.Bytes(`value`)
)

// AddHandler adds debug handlers to router, allows to add more middleware
// for example for access control
func AddHandlers(r *router.Group, prefix string, middleware ...router.MiddlewareFunc) {

	// clear state entries by key prefixs
	r.Invoke(
		prefix+InvokeStateCleanFunc,
		InvokeStateClean,
		append([]router.MiddlewareFunc{PrefixesParam}, middleware...)...)

	// query keys by prefix
	r.Query(
		prefix+QueryStateKeysFunc,
		QueryKeysList,
		append([]router.MiddlewareFunc{PrefixParam}, middleware...)...)

	// query value by key
	r.Query(
		prefix+QueryStateGetFunc,
		QueryStateGet,
		append([]router.MiddlewareFunc{KeyParam}, middleware...)...)

	r.Invoke(
		prefix+InvokeStatePutFunc,
		InvokeStatePut,
		append([]router.MiddlewareFunc{KeyParam, ValueParam}, middleware...)...)

	r.Invoke(
		prefix+InvokeStateDeleteFunc,
		InvokeStateDelete,
		append([]router.MiddlewareFunc{KeyParam}, middleware...)...)
}

// InvokeStateClean delete entries from state, prefix []string contains key prefixes or whole key
func InvokeStateClean(c router.Context) (interface{}, error) {
	return DeleteStateByPrefixes(c.State(), c.Param(`prefixes`).([]string))
}

// InvokeValueByKeyPut router handler puts value in chaincode state with composite key,
// created with key parts ([]string)
func InvokeStatePut(c router.Context) (interface{}, error) {
	key, err := state.KeyToString(c.Stub(), c.Param(`key`).([]string))
	if err != nil {
		return nil, errors.Wrap(err, `unable to create key`)
	}
	return nil, c.Stub().PutState(key, c.ParamBytes(`value`))
}

// QueryKeysList router handler returns string slice with keys by prefix (object type)
func QueryKeysList(c router.Context) (interface{}, error) {
	return c.State().Keys(c.ParamString(`prefix`))
}

// QueryStateGet router handler returns state entry by key ([]string)
func QueryStateGet(c router.Context) (interface{}, error) {
	key, err := state.KeyToString(c.Stub(), c.Param(`key`).([]string))
	if err != nil {
		return nil, errors.Wrap(err, `unable to create key`)
	}
	return c.Stub().GetState(key)
}

// QueryStateGet router handler delete state entry by key ([]string)
func InvokeStateDelete(c router.Context) (interface{}, error) {
	key, err := state.KeyToString(c.Stub(), c.Param(`key`).([]string))
	if err != nil {
		return nil, errors.Wrap(err, `unable to create key`)
	}
	return nil, c.Stub().DelState(key)
}
