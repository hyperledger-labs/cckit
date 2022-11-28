package envelope_test

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	e "github.com/hyperledger-labs/cckit/extensions/envelope"
	"github.com/hyperledger-labs/cckit/extensions/envelope/testdata"
	identitytestdata "github.com/hyperledger-labs/cckit/identity/testdata"
	"github.com/hyperledger-labs/cckit/serialize"

	testcc "github.com/hyperledger-labs/cckit/testing"
)

func TestEnvelop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Envelop suite")
}

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
			publicKey, privateKey, err := e.CreateKeys()
			Expect(err).NotTo(HaveOccurred())
			Expect(len(publicKey)).To(Equal(32))
			Expect(len(privateKey)).To(Equal(64))
		})

		It("Allow to create nonces", func() {
			nonce1 := e.CreateNonce()
			nonce2 := e.CreateNonce()

			Expect(nonce1).NotTo(BeEmpty())
			Expect(nonce2).NotTo(BeEmpty())
			// todo: test nonces equivalence
		})

		It("Allow to create signature", func() {
			_, privateKey, _ := e.CreateKeys()
			_, sig := e.CreateSig(payload, e.CreateNonce(), privateKey)
			Expect(len(sig)).To(Equal(64))
		})

		It("Allow to check valid signature", func() {
			nonce := e.CreateNonce()
			publicKey, privateKey, _ := e.CreateKeys()
			_, sig := e.CreateSig(payload, nonce, privateKey)
			err := e.CheckSig(payload, nonce, publicKey, sig)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Disallow to check signature with invalid payload", func() {
			nonce := e.CreateNonce()
			publicKey, privateKey, _ := e.CreateKeys()
			_, sig := e.CreateSig(payload, nonce, privateKey)
			invalidPayload := []byte("invalid payload")
			err := e.CheckSig(invalidPayload, nonce, publicKey, sig)
			Expect(err).Should(MatchError(e.ErrSignatureCheckFailed))
		})

	})

	Describe("Handle base64 envelop", func() {

		It("Allow to parse base64 envelop", func() {
			publicKey, privateKey, _ := e.CreateKeys()
			nonce := e.CreateNonce()
			hashToSign := e.Hash(payload, nonce)
			_, sig := e.CreateSig(payload, nonce, privateKey)
			envelope := &e.Envelope{
				PublicKey:       publicKey,
				Signature:       sig,
				Nonce:           nonce,
				HashToSign:      hashToSign[:],
				HashFunc:        "SHA3-256",
				Deadline:        testcc.MustProtoTimestamp(time.Now().AddDate(0, 2, 0)),
				DomainSeparator: []byte("DomainSeparator"),
			}
			jj, _ := json.Marshal(envelope)
			b64 := base64.StdEncoding.EncodeToString(jj)
			bb, err := e.DecodeEnvelope([]byte(b64))
			Expect(err).NotTo(HaveOccurred())
			Expect(bb).To(Equal([]byte(jj)))
		})

	})

	Describe("Signature verification", func() {

		var (
			serializedEnvelope []byte
			serializer         = serialize.DefaultSerializer
		)

		BeforeEach(func() {
			publicKey, privateKey, _ := e.CreateKeys()
			nonce := e.CreateNonce()
			hashToSign := e.Hash(payload, nonce)
			_, sig := e.CreateSig(payload, nonce, privateKey)
			envelope := &e.Envelope{
				PublicKey:       publicKey,
				Signature:       sig,
				Nonce:           nonce,
				HashToSign:      hashToSign[:],
				HashFunc:        "SHA3-256",
				Deadline:        testcc.MustProtoTimestamp(time.Now().AddDate(0, 2, 0)),
				DomainSeparator: []byte("DomainSeparator"),
			}
			serializedEnvelope, _ = serializer.ToBytesFrom(envelope)
		})

		It("Allow to verify valid signature", func() {
			envelopCC = testcc.NewMockStub(
				`envelop chaincode mock`,
				testdata.NewEnvelopCC(chaincodeName, channelName))

			resp := envelopCC.Invoke("invokeWithEnvelop", payload, serializedEnvelope)
			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Disallow to verify signature with invalid payload", func() {
			envelopCC = testcc.NewMockStub(
				`envelop chaincode mock`,
				testdata.NewEnvelopCC(chaincodeName, channelName))

			invalidPayload := []byte("invalid payload")
			resp := envelopCC.Invoke("invokeWithEnvelop", invalidPayload, serializedEnvelope)
			Expect(resp.Status).To(BeNumerically("==", 500))
		})

	})

	Describe("Nonce verification (replay attack)", func() {
		It("Disallow to execute tx with the same parameters (nonce, payload, pubkey)", func() {
			envelopCC = testcc.NewMockStub(
				`envelop chaincode mock`,
				testdata.NewEnvelopCC(chaincodeName, channelName))

			publicKey, privateKey, _ := e.CreateKeys()
			nonce := "thesamenonce"
			hashToSign := e.Hash(payload, nonce)
			_, sig := e.CreateSig(payload, nonce, privateKey)
			envelope := &e.Envelope{
				PublicKey:       publicKey,
				Signature:       sig,
				Nonce:           nonce,
				HashToSign:      hashToSign[:],
				HashFunc:        "SHA3-256",
				Deadline:        testcc.MustProtoTimestamp(time.Now().AddDate(0, 2, 0)),
				DomainSeparator: []byte("DomainSeparator"),
			}
			serializer := serialize.DefaultSerializer
			serializedEnvelope, _ := serializer.ToBytesFrom(envelope)

			resp := envelopCC.Invoke("invokeWithEnvelop", payload, serializedEnvelope)
			Expect(resp.Status).To(BeNumerically("==", 200))

			resp = envelopCC.Invoke("invokeWithEnvelop", payload, serializedEnvelope)
			Expect(errors.New(resp.Message)).To(MatchError(e.ErrTxAlreadyExecuted))
		})
	})

})