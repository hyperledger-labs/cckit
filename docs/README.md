# Getting started

## Hello world

Simplest chaincode with CCKit looks like

```golang
package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"

	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
)

func main() {
	// create cckit router
	r := router.New(`hello world`)

	// define chaincode method Hello with 1 argument `name`
	r.Query(`hello`, func(c router.Context) (interface{}, error) {
		return `Hello ` + c.ParamString(`name`), nil
	}, param.String(`name`))

	if err := shim.Start(router.NewChaincode(r)); err != nil {
		fmt.Printf("Error starting hello world chaincode: %s", err)
	}
}
```

## Recommended way with gRPC / protobuf definition and code generation

Using core CCKit components like 

* [Router](../router)
* [State mapper](../state)
* [Serializer](../state)

CCKit allows to generate chaincode interfaces and chaincode invocation layer.

Read more https://medium.com/coinmonks/service-oriented-hyperledger-fabric-application-development-32e66f578f9a

## Using extensions

### Encrypt chaincode state

Using [Encryption extension](../extensions/encryption)

### Store information about chaincode "owner"


Using [Owner extension](../extensions/owner)

### Create ERC-20 token

Using [Token extension](../extensions/token)
