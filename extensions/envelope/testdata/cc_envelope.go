package testdata

import (
	"github.com/hyperledger-labs/cckit/extensions/envelope"
	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
	"github.com/hyperledger-labs/cckit/serialize"
)

type EnvelopCC struct {
}

func NewEnvelopCC(chaincodeName string) *router.Chaincode {
	r := router.New(chaincodeName, serialize.PreferJSONSerializer).Use(envelope.Verify())

	r.Invoke("invokeWithEnvelope", func(c router.Context) (interface{}, error) {
		return nil, nil
	}, param.String("payload"), param.Bytes("envelope"))
	r.Query("queryWithoutEnvelope", func(c router.Context) (interface{}, error) {
		return nil, nil
	}, param.String("payload"))

	return router.NewChaincode(r)
}
