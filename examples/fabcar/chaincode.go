package fabcar

import (
	"github.com/hyperledger-labs/cckit/extensions/owner"
	"github.com/hyperledger-labs/cckit/router"
)

const ChaincodeName = `fabcar`

func New() (*router.Chaincode, error) {

	r := router.New(ChaincodeName)

	r.Init(ChaincodeInitFunc())

	if err := RegisterFabCarServiceChaincode(r, &FabCarService{}); err != nil {
		return nil, err
	}

	return router.NewChaincode(r), nil
}

func ChaincodeInitFunc() func(router.Context) (interface{}, error) {
	return func(ctx router.Context) (interface{}, error) {
		return owner.SetFromCreator(ctx)
	}
}
