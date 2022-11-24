package envelop

import (
	"context"
	"encoding/base64"
	"fmt"

	cckit_gateway "github.com/hyperledger-labs/cckit/gateway"
	"github.com/hyperledger-labs/cckit/serialize"
	"google.golang.org/grpc/metadata"
)

// custom invoker that can parse header with envelop
type ChaincodeInstanceServiceEnvelopInvoker struct {
	cckit_gateway.ChaincodeInstanceInvoker
	serializer *serialize.GenericSerializer
}

func NewChaincodeInstanceServiceEnvelopInvoker() *ChaincodeInstanceServiceEnvelopInvoker {
	return &ChaincodeInstanceServiceEnvelopInvoker{
		serializer: serialize.DefaultSerializer,
	}
}

// set default invoker to decorate calls
func (c *ChaincodeInstanceServiceEnvelopInvoker) DefaultInvoker(invoker cckit_gateway.ChaincodeInstanceInvoker) {
	c.ChaincodeInstanceInvoker = invoker
}

// stub to query default handler
func (c *ChaincodeInstanceServiceEnvelopInvoker) Query(
	ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {

	return c.ChaincodeInstanceInvoker.Query(ctx, fn, args, target)
}

// custom invoker for parsing envelop in headers
func (c *ChaincodeInstanceServiceEnvelopInvoker) Invoke(
	ctx context.Context, fn string, args []interface{}, target interface{}) (interface{}, error) {

	// get envelop with signature from header and add as second arg
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if v, ok := md["x-envelop"]; ok {
			env, err := ParseEnvelop([]byte(v[0]))
			if err != nil {
				return nil, fmt.Errorf(`invoke: %w`, err)
			}
			args = append(args, env)
		}
	}

	return c.ChaincodeInstanceInvoker.Invoke(ctx, fn, args, target)
}

// parse base64 envelop
func ParseEnvelop(base64Env []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(base64Env)))
	n, err := base64.StdEncoding.Decode(dst, base64Env)
	if err != nil {
		return nil, ErrDecodeEnvelopFailed
	}
	return dst[:n], nil
}
