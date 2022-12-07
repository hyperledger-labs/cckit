package testdata

import (
	"github.com/hyperledger-labs/cckit/extensions/envelope"
	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
)

type EnvelopCC struct {
}

func NewEnvelopCC(chaincodeName, channelName string) *router.Chaincode {
	// r := router.New(chaincodeName).Pre(envelope.PreVerify)
	r := router.New(chaincodeName).Use(envelope.Verify())

	r.Invoke("invokeWithEnvelope", func(c router.Context) (interface{}, error) {
		return nil, nil
	}, param.String("payload"), param.Bytes("envelope"))
	r.Query("queryWithoutEnvelope", func(c router.Context) (interface{}, error) {
		return nil, nil
	}, param.String("payload"))

	return router.NewChaincode(r)
}
