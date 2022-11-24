package testdata

import (
	"github.com/hyperledger-labs/cckit/extensions/envelop"
	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
)

type EnvelopCC struct {
}

func NewEnvelopCC(chaincodeName, channelName string) *router.Chaincode {
	r := router.New(chaincodeName).Pre(envelop.Verify)

	r.Invoke("invokeWithEnvelop", func(c router.Context) (interface{}, error) {
		return nil, nil
	}, param.String("payload"), param.String("env"))

	return router.NewChaincode(r)
}
