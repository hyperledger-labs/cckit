package envelop_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger-labs/cckit/extensions/envelop"
	"github.com/hyperledger-labs/cckit/extensions/envelop/testdata"
	identitytestdata "github.com/hyperledger-labs/cckit/identity/testdata"
	"github.com/hyperledger-labs/cckit/serialize"

	testcc "github.com/hyperledger-labs/cckit/testing"
)

func TestEnvelop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Envelop suite")
}

var ()

var (
	Owner = identitytestdata.Certificates[0].MustIdentity(`SOME_MSP`)

	envelopCC *testcc.MockStub

	chaincodeName = "envelop-chaincode"
	channelName   = "payment-channel"

	payload = []byte("signed payload")
)

var _ = Describe(`Envelop`, func() {

	Describe("Signature methods", func() {

		It("Allow to create keys", func() {
			publicKey, privateKey, err := envelop.CreateKeys()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(publicKey)).To(Equal(32))
			Expect(len(privateKey)).To(Equal(64))
		})

		It("Allow to create nonces", func() {
			nonce1 := envelop.CreateNonce()
			nonce2 := envelop.CreateNonce()

			Expect(nonce1).NotTo(BeEmpty())
			Expect(nonce2).NotTo(BeEmpty())
			// todo: test nonces equivalence
		})

		It("Allow to create signature", func() {
			_, privateKey, _ := envelop.CreateKeys()
			_, sig := envelop.CreateSig(payload, envelop.CreateNonce(), privateKey)
			Expect(len(sig)).To(Equal(64))
		})

		It("Allow to check valid signature", func() {
			nonce := envelop.CreateNonce()
			publicKey, privateKey, _ := envelop.CreateKeys()
			_, sig := envelop.CreateSig(payload, nonce, privateKey)
			err := envelop.CheckSig(payload, nonce, publicKey, sig)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Disallow to check signature with invalid payload", func() {
			nonce := envelop.CreateNonce()
			publicKey, privateKey, _ := envelop.CreateKeys()
			_, sig := envelop.CreateSig(payload, nonce, privateKey)
			invalidPayload := []byte("invalid payload")
			err := envelop.CheckSig(invalidPayload, nonce, publicKey, sig)
			Expect(err).Should(MatchError(envelop.ErrSignatureCheckFailed))
		})

	})

	Describe("Handle base64 envelop", func() {

		It("Allow to parse base64 envelop", func() {
			publicKey, privateKey, _ := envelop.CreateKeys()
			nonce := envelop.CreateNonce()
			hashToSign := envelop.Hash(payload, nonce)
			_, sig := envelop.CreateSig(payload, nonce, privateKey)
			env := &envelop.Envelop{
				PublicKey:       publicKey,
				Signature:       sig,
				Nonce:           nonce,
				HashToSign:      hashToSign[:],
				HashFunc:        "SHA3-256",
				Deadline:        testcc.MustProtoTimestamp(time.Now().AddDate(0, 2, 0)),
				DomainSeparator: []byte("DomainSeparator"),
			}
			jsonEnv, _ := json.Marshal(env)
			base64Env := base64.StdEncoding.EncodeToString(jsonEnv)
			bb, err := envelop.DecodeEnvelope([]byte(base64Env))
			Expect(err).NotTo(HaveOccurred())
			Expect(bb).To(Equal([]byte(jsonEnv)))
		})

	})

	Describe("Signature verification", func() {

		var (
			serializedEnv []byte
			serializer    = serialize.PreferJSONSerializer
		)

		BeforeEach(func() {
			publicKey, privateKey, _ := envelop.CreateKeys()
			nonce := envelop.CreateNonce()
			hashToSign := envelop.Hash(payload, nonce)
			_, sig := envelop.CreateSig(payload, nonce, privateKey)
			env := &envelop.Envelop{
				PublicKey:       publicKey,
				Signature:       sig,
				Nonce:           nonce,
				HashToSign:      hashToSign[:],
				HashFunc:        "SHA3-256",
				Deadline:        testcc.MustProtoTimestamp(time.Now().AddDate(0, 2, 0)),
				DomainSeparator: []byte("DomainSeparator"),
			}
			serializedEnv, _ = serializer.ToBytesFrom(env)
			// jsonEnv, _ = json.Marshal(env)
		})

		It("Allow to verify valid signature", func() {
			envelopCC = testcc.NewMockStub(
				`envelop chaincode mock`,
				testdata.NewEnvelopCC(chaincodeName, channelName))

			resp := envelopCC.Invoke("invokeWithEnvelop", payload, serializedEnv)
			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Disallow to verify signature with invalid payload", func() {
			envelopCC = testcc.NewMockStub(
				`envelop chaincode mock`,
				testdata.NewEnvelopCC(chaincodeName, channelName))

			invalidPayload := []byte("invalid payload")
			resp := envelopCC.Invoke("invokeWithEnvelop", invalidPayload, serializedEnv)
			Expect(resp.Status).To(BeNumerically("==", 500))
		})

	})

})
