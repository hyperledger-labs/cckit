# Hyperledger Fabric chaincode kit (CCKit)

[![Go Report Card](https://goreportcard.com/badge/github.com/hyperledger-labs/cckit)](https://goreportcard.com/report/github.com/hyperledger-labs/cckit)
[![Coverage Status](https://coveralls.io/repos/github/hyperledger-labs/cckit/badge.svg?branch=master)](https://coveralls.io/github/hyperledger-labs/cckit?branch=master)

## Overview

A [smart contract](https://hyperledger-fabric.readthedocs.io/en/latest/glossary.html#smart-contract) is code, 
invoked by a client application external to the blockchain network – that manages access and modifications to a set of
key-value pairs in the World State.  In Hyperledger Fabric, smart contracts are referred to as chaincode.

**CCKit** is a **programming toolkit** for

* developing and testing Hyperledger Fabric golang chaincodes 
* generating access layer for query / invoke chaincodes and chaincode event streaming

It enhances the development experience while providing developers components for creating more readable and secure
smart contracts. [Getting started]


### CCKit features

* [Chaincode method router](router) with invocation handlers and middleware capabilities
* [Chaincode state modeling](state) (like ORM for chaincode data) using protocol buffers or plain  golang struct with private data support 
* [Chaincode state serializing](serialize) customization
* Designing chaincode in [gRPC service notation](gateway) with code generation of chaincode SDK, gRPC and REST-API
* [MockStub testing](testing), allowing to immediately receive test results

### Extensions
* [Encryption](extensions/encryption) data on application level
* Implementation of UTXO and account balance storage model [token balance models](extensions/token)
* Chaincode method [access control](extensions/owner)
* Signing and verifying chaincode method payload with [envelope](extensions/envelope)


There are several chaincode "official" examples available:

* [Commercial paper](https://hyperledger-fabric.readthedocs.io/en/release-1.4/developapps/smartcontract.html) from official [Hyperledger Fabric documentation](https://hyperledger-fabric.readthedocs.io)
* [Blockchain insurance application](https://github.com/IBM/build-blockchain-insurance-app) (testing tutorial: how to [write tests for "insurance" chaincode](examples/insurance))

and others

**Main problems** with existing examples are:

* Working with chaincode state at very low level
* Lots of code duplication (JSON marshalling / unmarshalling, validation, access control, etc)
* Chaincode methods routing appeared only in HLF 1.4 and only in Node.Js chaincode
* Uncompleted testing tools (MockStub)

## Examples based on CCKit

* [Cars](examples/cars) - car registration chaincode, *simplest* example
* [Commercial paper, service-oriented approach](https://github.com/s7techlab/hyperledger-fabric-samples) - 
  recommended way to start new application. Code generation radically simplifies building on-chain and off-chain applications.
* [Commercial paper](examples/cpaper) - faithful reimplementation of the official example 
* [Commercial paper extended example](examples/cpaper_extended) - with protobuf chaincode state schema and other features
* [ERC-20](examples/erc20) - tokens smart contract, implementing ERC-20 interface
* [Cars private](examples/private_cars) - car registration chaincode with private data
* [Payment](examples/payment) - a few examples of chaincodes with encrypted state 

### Publications with usage examples

* [Service-oriented Hyperledger Fabric application development using gRPC definitions](https://medium.com/coinmonks/service-oriented-hyperledger-fabric-application-development-32e66f578f9a)
* [Hyperledger Fabric smart contract data model: protobuf to chaincode state mapping](https://medium.com/coinmonks/hyperledger-fabric-smart-contract-data-model-protobuf-to-chaincode-state-mapping-191cdcfa0b78)
* [Hyperledger Fabric chaincode test driven development (TDD) with unit testing](https://medium.com/coinmonks/test-driven-hyperledger-fabric-golang-chaincode-development-dbec4cb78049)
* [ERC20 token as Hyperledger Fabric Golang chaincode](https://medium.com/@viktornosov/erc20-token-as-hyperledger-fabric-golang-chaincode-d09dfd16a339)
* [CCKit: Routing and middleware for Hyperledger Fabric Golang chaincode](https://medium.com/@viktornosov/routing-and-middleware-for-developing-hyperledger-fabric-chaincode-written-in-go-90913951bf08)
* [Developing and testing Hyperledger Fabric smart contracts](https://habr.com/post/426705/) [RUS]


## Installation

CCKit requires Go 1.16+  

### Standalone
 
`git clone git@github.com:hyperledger-labs/cckit.git`

`go mod vendor`

### As dependency

`go get github.com/hyperledger-labs/cckit`

