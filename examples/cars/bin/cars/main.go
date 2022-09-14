package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/hyperledger-labs/cckit/examples/cars"
)

func main() {
	cc := cars.New()
	if err := shim.Start(cc); err != nil {
		fmt.Printf("Error starting Cars chaincode: %s", err)
	}
}
