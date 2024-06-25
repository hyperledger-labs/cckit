package envelope_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger-labs/cckit/extensions/envelope"
)

const (
	Ed25519PublicKeyLen  = 32
	Ed25519PrivateKeyLen = 64
	Ed25519SignatureLen  = 64
)

var _ = Describe(`Ed25519 crypto`, func() {

	crypto := envelope.NewEd25519()

	It("Allow to create keys", func() {
		publicKey, privateKey, err := crypto.GenerateKey()
		Expect(err).NotTo(HaveOccurred())
		Expect(len(publicKey)).To(Equal(Ed25519PublicKeyLen))
		Expect(len(privateKey)).To(Equal(Ed25519PrivateKeyLen))
	})

	It("Allow to create signature", func() {
		_, privateKey, _ := crypto.GenerateKey()
		sig, err := crypto.Sign(privateKey, []byte(`anything`))
		Expect(err).NotTo(HaveOccurred())
		Expect(len(sig)).To(Equal(Ed25519SignatureLen))
	})
})
