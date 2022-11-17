package router_test

import (
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger-labs/cckit/router"
	testcc "github.com/hyperledger-labs/cckit/testing"
)

func TestRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Router suite")
}

func New() *router.Chaincode {
	r := router.New(`router`).
		Init(router.EmptyContextHandler).
		Invoke(`empty`, func(c router.Context) (interface{}, error) {
			return nil, nil
		})

	return router.NewChaincode(r)
}

var cc *testcc.MockStub

var _ = Describe(`Router`, func() {

	BeforeSuite(func() {
		cc = testcc.NewMockStub(`Router`, New())
	})

	It(`Allow empty response`, func() {

		response := cc.Invoke(`empty`)
		Expect(response.Status).To(Equal(int32(shim.OK)))
		Expect(response.Payload).To(BeEmpty())
		Expect(response.Message).To(BeEmpty())

	})

})
