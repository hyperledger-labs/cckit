package gateway

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/hyperledger-labs/cckit/serialize"
)

func invokerArgs(fn string, args []interface{}, toBytesConverter serialize.ToBytesConverter) ([][]byte, error) {
	argsBytes, err := serialize.ArgsToBytes(args, toBytesConverter)
	if err != nil {
		return nil, fmt.Errorf(`invoker args: %w`, err)
	}

	return append([][]byte{[]byte(fn)}, argsBytes...), nil
}

func ccInput(ctx context.Context, fn string, args []interface{}, toBytesConverter serialize.ToBytesConverter) (*ChaincodeInput, error) {
	argsBytes, err := invokerArgs(fn, args, toBytesConverter)
	if err != nil {
		return nil, fmt.Errorf(`input: %w`, err)
	}
	ccInput := &ChaincodeInput{
		Args: argsBytes,
	}

	if ccInput.Transient, err = TransientFromContext(ctx); err != nil {
		return nil, err
	}

	return ccInput, nil
}

func ccOutput(response *peer.Response, target interface{}, fromBytesConverter serialize.FromBytesConverter) (res interface{}, err error) {
	output, err := fromBytesConverter.FromBytesTo(response.Payload, target)
	if err != nil {
		return nil, fmt.Errorf(`output to=%s: %w`, reflect.TypeOf(target), err)
	}

	return output, nil
}
