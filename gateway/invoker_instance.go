package gateway

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/cckit/serialize"
)

// ChaincodeInvoker used in generated service gateway code
type (
	ChaincodeInstanceInvoker interface {
		Query(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error)
		Invoke(ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error)
	}

	ChaincodeInstanceServiceInvoker struct {
		ChaincodeInstance ChaincodeInstanceServiceServer
		Serializer        serialize.Serializer
	}
)

// NewChaincodeInstanceServiceInvoker
func NewChaincodeInstanceServiceInvoker(
	ccInstance ChaincodeInstanceServiceServer, opts ...InvokerOpt) *ChaincodeInstanceServiceInvoker {
	c := &ChaincodeInstanceServiceInvoker{
		ChaincodeInstance: ccInstance,
		Serializer:        serialize.DefaultSerializer, // set default serializer as proto
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *ChaincodeInstanceServiceInvoker) Query(
	ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {

	ccInput, err := ccInput(ctx, fn, args, c.Serializer)
	if err != nil {
		return nil, fmt.Errorf(`query: %w`, err)
	}

	res, err := c.ChaincodeInstance.Query(ctx, &ChaincodeInstanceQueryRequest{
		Input: ccInput,
	})
	if err != nil {
		return nil, err
	}

	return ccOutput(res, target, c.Serializer)
}

func (c *ChaincodeInstanceServiceInvoker) Invoke(
	ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {

	ccInput, err := ccInput(ctx, fn, args, c.Serializer)
	if err != nil {
		return nil, fmt.Errorf(`invoke: %w`, err)
	}

	res, err := c.ChaincodeInstance.Invoke(ctx, &ChaincodeInstanceInvokeRequest{
		Input: ccInput,
	})
	if err != nil {
		return nil, err
	}

	return ccOutput(res, target, c.Serializer)
}
