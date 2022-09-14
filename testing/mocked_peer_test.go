package testing_test

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger-labs/cckit/examples/cpaper_asservice"
	cpservice "github.com/hyperledger-labs/cckit/examples/cpaper_asservice"
	"github.com/hyperledger-labs/cckit/gateway"
	idtestdata "github.com/hyperledger-labs/cckit/identity/testdata"
	testcc "github.com/hyperledger-labs/cckit/testing"
)

var _ = Describe(`Service`, func() {

	const (
		ChaincodeName = `commercial_paper`
	)

	var (
		peer          *testcc.MockedPeerDecorator
		cPaperGateway *cpservice.CPaperServiceGateway

		ctx = gateway.ContextWithSigner(
			context.Background(),
			idtestdata.Certificates[0].MustIdentity(idtestdata.DefaultMSP),
		)
	)

	It("Init", func() {
		ccImpl, err := cpaper_asservice.NewCC()
		Expect(err).NotTo(HaveOccurred())

		// peer imitation
		peer = testcc.NewPeerDecorator(testcc.NewPeer().WithChannel(Channel, testcc.NewMockStub(ChaincodeName, ccImpl)))

		// "sdk" for deal with cpaper chaincode
		cPaperGateway = cpservice.NewCPaperServiceGateway(peer, Channel, ChaincodeName)
	})

	It("Default invoker", func() {
		_, err := cPaperGateway.List(ctx, &empty.Empty{})
		Expect(err).NotTo(HaveOccurred())
	})

	It("Allow to imitate error while access to peer", func() {
		peer.FailChaincode(ChaincodeName)

		_, err := cPaperGateway.List(ctx, &empty.Empty{})
		Expect(err).To(HaveOccurred())
	})
})
