package encryption

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/pkg/errors"

	"github.com/hyperledger-labs/cckit/serialize"
	"github.com/hyperledger-labs/cckit/state"
)

// InvokeChaincode decrypts received payload
func InvokeChaincode(
	stub shim.ChaincodeStubInterface, encKey []byte, chaincodeName string,
	args []interface{}, channel string, target interface{}, serializer serialize.Serializer) (interface{}, error) {

	// args are not encrypted cause we cannot pass encryption key in transient map while invoking cc from cc
	// thus target cc cannot decrypt args
	aa, err := serialize.ArgsToBytes(args, serializer)
	if err != nil {
		return nil, errors.Wrap(err, `encrypt args`)
	}

	response := stub.InvokeChaincode(chaincodeName, aa, channel)
	if response.Status != shim.OK {
		return nil, errors.New(response.Message)
	}

	if len(response.Payload) == 0 {
		return nil, state.ErrEmptyChaincodeResponsePayload
	}

	decrypted, err := Decrypt(encKey, response.Payload)
	if err != nil {
		return nil, errors.Wrap(err, `decrypt payload`)
	}
	return serializer.FromBytesTo(decrypted, target)
}
