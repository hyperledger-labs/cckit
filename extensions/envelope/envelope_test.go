package envelope_test

import (
	"encoding/base64"
	"encoding/hex"
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

	chaincode    = "envelope-chaincode"
	channel      = "envelope-channel"
	methodInvoke = "invokeWithEnvelope"
	methodQuery  = "queryWithoutEnvelope"
	deadline     = timestamppb.New(time.Now().AddDate(0, 0, 2))

	payload = []byte(`{"symbol":"GLD","decimals":"8","name":"Gold digital asset","type":"DM","underlying_asset":"gold","issuer_id":"GLDINC"}`)
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

		It("Allow to verify valid signature from the envelope in base64 format with deadline", func() {
			b64Envelope := "eyJoYXNoX2Z1bmMiOiJTSEEzLTI1NiIsImhhc2hfdG9fc2lnbiI6ImE3MTdlYTY2MzM5ZmM1MzVmNWY4NzIxODkyNmFjYzZkODJlOTk1NGM3Y2YwNGRlZjdkZGI1YzA3NmRmNWY4OTUiLCJub25jZSI6IjE2NzQ3MzU1NTkyMDEiLCJjaGFubmVsIjoiZW52ZWxvcGUtY2hhbm5lbCIsIm1ldGhvZCI6Imludm9rZVdpdGhFbnZlbG9wZSIsImNoYWluY29kZSI6ImVudmVsb3BlLWNoYWluY29kZSIsImRlYWRsaW5lIjoiMjAyMy0wMS0yN1QxMjoxOToxOS4yMDFaIiwicHVibGljX2tleSI6IjFmMjQyOGY3M2JkMWM1ZDM0MGExYTZlNDMxMzk1Y2Q0NDk0MzRhZjBlMjg2NTFiYzdlYjQyOGIyYmNmZWRjNmUiLCJzaWduYXR1cmUiOiIyYWYzYjVjNmE2Yjg1MTk4MjllMWMwOGY3ZjYxNzEzODQyYmVjM2I2Y2QxN2EzZTkyM2ZhMDFkNTVkNGY2OWU2ZGI5YzkyOTdiNGRmZDcxNDRlOWQ1YWNjZGRkMDA2M2RiZDc5M2U1ZDIxYjhmMTllY2QyYjdhMmJhOGU3YjMwZiJ9"
			decodedEnvelope, err := e.DecodeEnvelope([]byte(b64Envelope))
			Expect(err).NotTo(HaveOccurred())

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(chaincode)).WithChannel(channel)
			resp := envelopCC.Invoke(methodInvoke, payload, decodedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Allow to verify valid signature from the envelope in base64 format without deadline", func() {
			b64Envelope := "eyJoYXNoX2Z1bmMiOiJTSEEyNTYiLCJoYXNoX3RvX3NpZ24iOiJiM2YyODQ4YTM4NThhNWVlY2NiOTI2MTI1NzA2ZjdhYjk4MmI0OWFlZjI3YmRlNDI3M2QzYWVkMDUwMTVmMzBjIiwibm9uY2UiOiIxNjc0NzM4ODM2NzExIiwiY2hhbm5lbCI6ImVudmVsb3BlLWNoYW5uZWwiLCJtZXRob2QiOiJpbnZva2VXaXRoRW52ZWxvcGUiLCJjaGFpbmNvZGUiOiJlbnZlbG9wZS1jaGFpbmNvZGUiLCJkZWFkbGluZSI6IjE5NzAtMDEtMDFUMDA6MDA6MDAuMDAwWiIsInB1YmxpY19rZXkiOiIyNDdkYmRiZmYzZmE1MGYwZWFiZTU2NDE4ZWMxMTA2ODk4ZTQ0ZDc1ZTc2MTYyNzcyZDA4ZWJmZDVkY2IxYWFlIiwic2lnbmF0dXJlIjoiNjRmNTI5MDU2NzkyOTRiOGRjMzNhYjRjZmUyZjE0NmFjYmMxYmE2MmVjZTcxNzAyZmYxZjgwNDRiNzllZmRiODcwNWE4NjUyYmU5NGQ0YzViZTBiZDhiZWU1YWUzNGUzMmI4NDVmNWMzZWFlODc3MGE3MjBjNjkyYWIzODZiMGEifQ=="
			decodedEnvelope, err := e.DecodeEnvelope([]byte(b64Envelope))
			Expect(err).NotTo(HaveOccurred())

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(chaincode)).WithChannel(channel)
			resp := envelopCC.Invoke(methodInvoke, payload, decodedEnvelope)

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
			hashToSign := e.Hash(payload, nonce, channel, chaincode, methodInvoke, deadline.AsTime().Format(e.TimeLayout), []byte(publicKey))
			_, sig := e.CreateSig(payload, nonce, channel, chaincode, methodInvoke, deadline.AsTime().Format(e.TimeLayout), privateKey)
			envelope := &e.Envelope{
				PublicKey:  hex.EncodeToString([]byte(publicKey)),
				Signature:  hex.EncodeToString(sig),
				Nonce:      nonce,
				HashToSign: hex.EncodeToString(hashToSign[:]),
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
		PublicKey: hex.EncodeToString([]byte(publicKey)),
		Nonce:     nonce,
		HashFunc:  "SHA256",
		Channel:   channel,
		Chaincode: chaincode,
		Method:    method,
	}
	var formatDeadline string
	if len(deadline) > 0 {
		envelope.Deadline = deadline[0]
		formatDeadline = envelope.Deadline.AsTime().Format(e.TimeLayout)
	}
	hashToSign := e.Hash(payload, nonce, channel, chaincode, method, formatDeadline, publicKey)
	envelope.HashToSign = hex.EncodeToString(hashToSign[:])

	_, sig := e.CreateSig(payload, nonce, channel, chaincode, method, formatDeadline, privateKey)
	envelope.Signature = hex.EncodeToString(sig)

	serializedEnvelope, _ := serialize.PreferJSONSerializer.ToBytesFrom(envelope)
	return serializedEnvelope, envelope
}
