package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
)

func main() {
	r := router.New(`hello world`)

	// define chaincode method Hello with 1 argument `name`
	r.Query(`hello`, func(c router.Context) (interface{}, error) {
		return `Hello ` + c.ParamString(`name`), nil
	}, param.String(`name`))

	if err := shim.Start(router.NewChaincode(r)); err != nil {
		fmt.Printf("Error starting hello world chaincode: %s", err)
	}
}
