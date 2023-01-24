package envelope_test

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

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

	chaincode    = "envelop-chaincode"
	channel      = "inv-channel"
	methodInvoke = "invokeWithEnvelope"
	methodQuery  = "queryWithoutEnvelope"
	deadline     = timestamppb.New(time.Now().AddDate(0, 0, 2))

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
			_, sig := e.CreateSig(payload, e.CreateNonce(), channel, chaincode, methodInvoke, deadline.String(), privateKey)
			Expect(len(sig)).To(Equal(64))
		})

		It("Allow to check valid signature", func() {
			nonce := e.CreateNonce()
			publicKey, privateKey, _ := e.CreateKeys()
			_, sig := e.CreateSig(payload, nonce, channel, chaincode, methodInvoke, deadline.String(), privateKey)
			err := e.CheckSig(payload, nonce, channel, chaincode, methodInvoke, deadline.String(), publicKey, sig)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Disallow to check signature with invalid payload", func() {
			nonce := e.CreateNonce()
			publicKey, privateKey, _ := e.CreateKeys()
			_, sig := e.CreateSig(payload, nonce, channel, chaincode, methodInvoke, deadline.String(), privateKey)
			invalidPayload := []byte("invalid payload")
			err := e.CheckSig(invalidPayload, nonce, channel, chaincode, methodInvoke, deadline.String(), publicKey, sig)
			Expect(err).Should(MatchError(e.ErrSignatureCheckFailed))
		})

	})

	Describe("Handle base64 envelop", func() {

		It("Allow to parse base64 envelop", func() {
			_, envelope := createEnvelope(payload, channel, chaincode, methodInvoke, deadline)
			jj, _ := json.Marshal(envelope)
			b64 := base64.StdEncoding.EncodeToString(jj)
			bb, err := e.DecodeEnvelope([]byte(b64))
			Expect(err).NotTo(HaveOccurred())
			Expect(bb).To(Equal([]byte(jj)))
		})

	})

	Describe("Signature verification", func() {

		It("Allow to verify valid signature", func() {
			serializedEnvelope, _ := createEnvelope(payload, channel, chaincode, methodInvoke, deadline)

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(chaincode)).WithChannel(channel)
			resp := envelopCC.Invoke(methodInvoke, payload, serializedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Allow to verify valid signature without deadline", func() {
			serializedEnvelope, _ := createEnvelope(payload, channel, chaincode, methodInvoke)

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(chaincode)).WithChannel(channel)
			resp := envelopCC.Invoke(methodInvoke, payload, serializedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Disallow to verify signature with invalid payload", func() {
			serializedEnvelope, _ := createEnvelope(payload, channel, chaincode, methodInvoke, deadline)

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(chaincode)).WithChannel(channel)
			invalidPayload := []byte("invalid payload")
			resp := envelopCC.Invoke(methodInvoke, invalidPayload, serializedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 500))
		})

		It("Disallow to verify signature with invalid method", func() {
			serializedEnvelope, _ := createEnvelope(payload, channel, chaincode, "invalid method", deadline)

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(chaincode)).WithChannel(channel)
			resp := envelopCC.Invoke(methodInvoke, payload, serializedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 500))
		})

		It("Disallow to verify signature with invalid channel", func() {
			serializedEnvelope, _ := createEnvelope(payload, "invalid channel", chaincode, methodInvoke, deadline)

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(chaincode)).WithChannel(channel)
			resp := envelopCC.Invoke(methodInvoke, payload, serializedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 500))
		})

		It("Don't check signature for query method", func() {
			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(chaincode)).WithChannel(channel)
			resp := envelopCC.Query(methodQuery, payload)

			Expect(resp.Status).To(BeNumerically("==", 200))
		})

	})

	Describe("Nonce verification (replay attack)", func() {
		It("Disallow to execute tx with the same parameters (nonce, payload, pubkey)", func() {
			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(chaincode)).WithChannel(channel)

			publicKey, privateKey, _ := e.CreateKeys()
			nonce := "thesamenonce"
			hashToSign := e.Hash(payload, nonce, channel, chaincode, methodInvoke, deadline.String(), []byte(publicKey))
			_, sig := e.CreateSig(payload, nonce, channel, chaincode, methodInvoke, deadline.String(), privateKey)
			envelope := &e.Envelope{
				PublicKey:  publicKey,
				Signature:  sig,
				Nonce:      nonce,
				HashToSign: hashToSign[:],
				HashFunc:   "SHA3-256",
				Deadline:   deadline,
				Channel:    channel,
				Chaincode:  chaincode,
				Method:     methodInvoke,
			}
			serializer := serialize.PreferJSONSerializer
			serializedEnvelope, _ := serializer.ToBytesFrom(envelope)

			resp := envelopCC.Invoke(methodInvoke, payload, serializedEnvelope)
			Expect(resp.Status).To(BeNumerically("==", 200))

			resp = envelopCC.Invoke(methodInvoke, payload, serializedEnvelope)
			Expect(errors.New(resp.Message)).To(MatchError(e.ErrTxAlreadyExecuted))
		})
	})

})

func createEnvelope(payload []byte, channel, chaincode, method string, deadline ...*timestamppb.Timestamp) ([]byte, *e.Envelope) {
	publicKey, privateKey, _ := e.CreateKeys()
	nonce := e.CreateNonce()
	envelope := &e.Envelope{
		PublicKey: publicKey,
		Nonce:     nonce,
		HashFunc:  "SHA3-256",
		Channel:   channel,
		Chaincode: chaincode,
		Method:    method,
	}
	if len(deadline) > 0 {
		envelope.Deadline = deadline[0]
	}
	hashToSign := e.Hash(payload, nonce, channel, chaincode, method, envelope.Deadline.String(), []byte(publicKey))
	_, sig := e.CreateSig(payload, nonce, channel, chaincode, method, envelope.Deadline.String(), privateKey)
	envelope.HashToSign = hashToSign[:]
	envelope.Signature = sig

	serializedEnvelope, _ := serialize.PreferJSONSerializer.ToBytesFrom(envelope)
	return serializedEnvelope, envelope
}
