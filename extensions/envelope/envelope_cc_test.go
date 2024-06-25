package envelope_test

import (
	"github.com/hyperledger-labs/cckit/extensions/envelope"
	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
	"github.com/hyperledger-labs/cckit/serialize"
	testcc "github.com/hyperledger-labs/cckit/testing"
)

type EnvelopCC struct {
}

const (
	chaincode    = "envelope-chaincode"
	channel      = "envelope-channel"
	methodInvoke = "invokeWithEnvelope"
	methodQuery  = "queryWithoutEnvelope"
)

func NewNewEnvelopCCMock(verifier envelope.Verifier) *testcc.MockStub {
	return testcc.NewMockStub(chaincode, NewEnvelopCC(verifier, chaincode)).WithChannel(channel)
}

func NewEnvelopCC(verifier envelope.Verifier, chaincodeName string) *router.Chaincode {
	r := router.New(chaincodeName, router.WithSerializer(serialize.PreferJSONSerializer)).Use(envelope.Verify(verifier))

	r.Invoke("invokeWithEnvelope", func(c router.Context) (interface{}, error) {
		return nil, nil
	}, param.String("payload"), param.Bytes("envelope"))
	r.Query("queryWithoutEnvelope", func(c router.Context) (interface{}, error) {
		return nil, nil
	}, param.String("payload"))

	return router.NewChaincode(r)
}
