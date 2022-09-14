package pinger

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger-labs/cckit/identity/testdata"
	"github.com/hyperledger-labs/cckit/router"
	testcc "github.com/hyperledger-labs/cckit/testing"
	expectcc "github.com/hyperledger-labs/cckit/testing/expect"
)

func TestPinger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pinger suite")
}

var (
	Someone = testdata.Certificates[1].MustIdentity(`SOME_MSP`)
)

func New() *router.Chaincode {
	r := router.New(`pingable`).
		Init(router.EmptyContextHandler).
		Invoke(FuncPing, Ping)
	return router.NewChaincode(r)
}

var _ = Describe(`Pinger`, func() {

	// Create chaincode mock
	cc := testcc.NewMockStub(`cars`, New())

	Describe("Pinger", func() {

		It("Allow anyone to invoke ping method", func() {
			//invoke chaincode method from authority actor
			pingInfo := expectcc.PayloadIs(cc.From(Someone).Invoke(FuncPing), &PingInfo{}).(*PingInfo)
			Expect(pingInfo.InvokerId).To(Equal(Someone.GetID()))
			Expect(pingInfo.InvokerCert).To(Equal(Someone.GetPEM()))
			Expect(pingInfo.EndorsingServerTime).To(Not(BeNil()))
			Expect(pingInfo.TxTime).To(Not(BeNil()))
		})
	})
})
