package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/hyperledger-labs/cckit/examples/erc20_utxo"
	"github.com/hyperledger-labs/cckit/examples/fabcar"
)

func main() {
	cc, err := erc20_utxo.NewChaincode()
	if err != nil {
		fmt.Printf("error creating %s chaincode: %s", fabcar.ChaincodeName, err)
		return
	}
	err = shim.Start(cc)
	if err != nil {
		fmt.Printf("Error starting ERC-20 chaincode: %s", err)
	}
}
