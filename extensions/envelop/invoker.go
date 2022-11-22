package envelop

import (
	"context"

	cckit_gateway "github.com/hyperledger-labs/cckit/gateway"
	"google.golang.org/grpc/metadata"
)

// custom invoker that can parse header with envelop
type ChaincodeInstanceServiceEnvelopInvoker struct {
	cckit_gateway.ChaincodeInstanceInvoker
}

// constructor that is passed to the cckit_gateway
func NewEnvelopInvoker(invoker cckit_gateway.ChaincodeInstanceInvoker) cckit_gateway.ChaincodeInstanceInvoker {
	c := &ChaincodeInstanceServiceEnvelopInvoker{
		ChaincodeInstanceInvoker: invoker,
	}

	return c
}

// stub to query default handler
func (c *ChaincodeInstanceServiceEnvelopInvoker) Query(
	ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {

	return c.ChaincodeInstanceInvoker.Query(ctx, fn, args, target)
}

// custom invoke to parse header with envelop
func (c *ChaincodeInstanceServiceEnvelopInvoker) Invoke(
	ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {

	// get envelop with signature from header and add as second arg
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if v, ok := md["x-envelop"]; ok {
			args = append(args, v[0])
		}
	}
	return c.ChaincodeInstanceInvoker.Invoke(ctx, fn, args, target)
}
