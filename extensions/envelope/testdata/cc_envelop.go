package testdata

import (
	"github.com/hyperledger-labs/cckit/extensions/envelope"
	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
)

type EnvelopCC struct {
}

func NewEnvelopCC(chaincodeName, channelName string) *router.Chaincode {
	r := router.New(chaincodeName).Pre(envelope.Verify)

	r.Invoke("invokeWithEnvelop", func(c router.Context) (interface{}, error) {
		return nil, nil
	}, param.String("payload"), param.Bytes("envelope"))

	return router.NewChaincode(r)
}
