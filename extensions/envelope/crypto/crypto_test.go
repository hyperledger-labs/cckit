package crypto_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger-labs/cckit/extensions/envelope/crypto"
)

const (
	Ed25519PublicKeyLen  = 32
	Ed25519PrivateKeyLen = 64
	Ed25519SignatureLen  = 64
)

func TestCrypto(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crypto suite")
}

var _ = Describe(`Ed25519 crypto`, func() {

	ed25519 := crypto.NewEd25519()

	It("Allow to create keys", func() {
		publicKey, privateKey, err := ed25519.GenerateKey()
		Expect(err).NotTo(HaveOccurred())
		Expect(len(publicKey)).To(Equal(Ed25519PublicKeyLen))
		Expect(len(privateKey)).To(Equal(Ed25519PrivateKeyLen))
	})

	It("Allow to create signature", func() {
		_, privateKey, _ := ed25519.GenerateKey()
		sig, err := ed25519.Sign(privateKey, []byte(`anything`))
		Expect(err).NotTo(HaveOccurred())
		Expect(len(sig)).To(Equal(Ed25519SignatureLen))
	})
})
