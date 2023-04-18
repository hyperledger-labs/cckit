package mapping_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	identitytestdata "github.com/hyperledger-labs/cckit/identity/testdata"
	"github.com/hyperledger-labs/cckit/serialize"
	"github.com/hyperledger-labs/cckit/state/mapping/testdata"
	"github.com/hyperledger-labs/cckit/state/mapping/testdata/schema"
	testcc "github.com/hyperledger-labs/cckit/testing"
	expectcc "github.com/hyperledger-labs/cckit/testing/expect"
)

func TestState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "State suite")
}

var (
	compositeIDCC, complexIDCC, sliceIDCC, indexesCC, configCC *testcc.MockStub

	Owner = identitytestdata.Certificates[0].MustIdentity(`SOME_MSP`)
)
var _ = Describe(`State mapping in chaincode`, func() {

	Describe("init chaincode", func() {
		compositeIDCC = testcc.NewMockStub(`proto`, testdata.NewCompositeIdCC(serialize.PreferJSONSerializer))
		compositeIDCC.Serializer = serialize.PreferJSONSerializer // need to set for correct invoking
		compositeIDCC.From(Owner).Init()
	})

	It("Allow to add data to chaincode state", func() {
		expectcc.ResponseOk(compositeIDCC.Invoke(testdata.CreateFunc, testdata.CreateEntityWithCompositeId[0]))
	})

	It("Allow to get entry list", func() {
		res := compositeIDCC.Query(testdata.ListFunc)
		Expect(string(res.Payload)[0:1]).To(Equal(`{`)) // json serialized
		entities := expectcc.JSONPayloadIs(res, &schema.EntityWithCompositeIdList{}).(*schema.EntityWithCompositeIdList)
		Expect(len(entities.Items)).To(Equal(1))
		Expect(entities.Items[0].Name).To(Equal(testdata.CreateEntityWithCompositeId[0].Name))
		Expect(entities.Items[0].Value).To(BeNumerically("==", testdata.CreateEntityWithCompositeId[0].Value))
	})
})
