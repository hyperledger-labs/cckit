package serialize_test

import (
	"testing"

	"github.com/hyperledger-labs/cckit/serialize"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "State suite")
}

var _ = Describe(`Generic serializer`, func() {

	serializer := serialize.DefaultSerializer

	It(`Bool`, func() {
		bTrue, err := serializer.ToBytesFrom(true)
		Expect(err).NotTo(HaveOccurred())
		Expect(bTrue).To(Equal([]byte(`true`)))

		bFalse, err := serializer.ToBytesFrom(false)
		Expect(err).NotTo(HaveOccurred())
		Expect(bFalse).To(Equal([]byte(`false`)))

		eTrue, err := serializer.FromBytesTo(bTrue, serialize.TypeBool)
		Expect(err).NotTo(HaveOccurred())
		Expect(eTrue.(bool)).To(Equal(true))

		eFalse, err := serializer.FromBytesTo(bFalse, serialize.TypeBool)
		Expect(err).NotTo(HaveOccurred())
		Expect(eFalse.(bool)).To(Equal(false))
	})

	It(`String`, func() {
		const MyStr = `my-string`
		bStr, err := serializer.ToBytesFrom(MyStr)
		Expect(err).NotTo(HaveOccurred())
		Expect(bStr).To(Equal([]byte(MyStr)))

		eStr, err := serializer.FromBytesTo(bStr, serialize.TypeString)
		Expect(err).NotTo(HaveOccurred())
		Expect(eStr.(string)).To(Equal(MyStr))
	})

	It(`Nil`, func() {
		bNil, err := serializer.ToBytesFrom(nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(bNil).To(Equal([]byte{}))
	})

})
