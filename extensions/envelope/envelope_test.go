package envelope_test

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/base58"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"

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

var _ = Describe(`Envelop with Ed25519`, func() {

	ed25519 := e.NewEd25519()
	signer := e.NewSigner(ed25519)

	Describe("Signature methods", func() {

		It("Allow to create keys", func() {
			publicKey, privateKey, err := ed25519.GenerateKey()
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
			_, privateKey, _ := ed25519.GenerateKey()
			sig, err := signer.Sign(payload, e.CreateNonce(), channel, chaincode, methodInvoke, deadline.String(), privateKey)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(sig)).To(Equal(64))
		})

		It("Allow to check valid signature", func() {
			nonce := e.CreateNonce()
			publicKey, privateKey, _ := ed25519.GenerateKey()
			sig, err := signer.Sign(payload, nonce, channel, chaincode, methodInvoke, deadline.String(), privateKey)
			Expect(err).NotTo(HaveOccurred())
			err = signer.CheckSignature(payload, nonce, channel, chaincode, methodInvoke, deadline.String(), publicKey, sig)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Disallow to check signature with invalid payload", func() {
			nonce := e.CreateNonce()
			publicKey, privateKey, _ := ed25519.GenerateKey()
			sig, _ := signer.Sign(payload, nonce, channel, chaincode, methodInvoke, deadline.String(), privateKey)
			invalidPayload := []byte("invalid payload")
			err := signer.CheckSignature(invalidPayload, nonce, channel, chaincode, methodInvoke, deadline.String(), publicKey, sig)
			Expect(err).Should(MatchError(e.ErrSignatureCheckFailed))
		})
	})

	Describe("Handle base64 envelop", func() {

		It("Allow to parse base64 envelop", func() {
			_, envelope := createEnvelope(signer, payload, channel, chaincode, methodInvoke, deadline)
			jj, _ := json.Marshal(envelope)
			b64 := base64.StdEncoding.EncodeToString(jj)
			bb, err := e.DecodeEnvelope([]byte(b64))
			Expect(err).NotTo(HaveOccurred())
			Expect(bb).To(Equal([]byte(jj)))
		})

	})

	Describe("Signature verification", func() {

		It("Allow to verify valid signature", func() {
			serializedEnvelope, _ := createEnvelope(signer, payload, channel, chaincode, methodInvoke, deadline)

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(signer, chaincode)).WithChannel(channel)
			resp := envelopCC.Invoke(methodInvoke, payload, serializedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Allow to verify valid signature without deadline", func() {
			serializedEnvelope, _ := createEnvelope(signer, payload, channel, chaincode, methodInvoke)

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(signer, chaincode)).WithChannel(channel)
			resp := envelopCC.Invoke(methodInvoke, payload, serializedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Allow to verify valid signature from the envelope in base64 format with deadline", func() {
			jsonEnvelope := `{"hash_func":"SHA256","hash_to_sign":"Bjd5D5qq3FFj1fSGt7QsHNqrfdDBzY7Cs5Wm45WZ19LE","nonce":"1675065805271","channel":"envelope-channel","method":"invokeWithEnvelope","chaincode":"envelope-chaincode","deadline":"2033-01-31T07:58:39.677Z","public_key":"EH9cLNxpg4FQmow9i1cF1b2vXkaJ17wUym7GSEeX6LQv","signature":"VG8P1ShwTGqVNhV8DVaZaVfxbZp7E8G9cpBbUPXFsGiwqLk1NJZFg2jt3ff1uJK92t2TbzkfyiL5ZTHsjYCdoQk"}`
			b64Envelope := base64.StdEncoding.EncodeToString([]byte(jsonEnvelope))
			decodedEnvelope, err := e.DecodeEnvelope([]byte(b64Envelope))
			Expect(err).NotTo(HaveOccurred())

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(signer, chaincode)).WithChannel(channel)
			resp := envelopCC.Invoke(methodInvoke, payload, decodedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Allow to verify valid signature from the envelope in base64 format without deadline", func() {
			jsonEnvelope := `{"hash_func":"SHA256","hash_to_sign":"H4HmQKUQJm2bxJvSDPpHaP5vYYGvL5dUqWEfTzNN3eeH","nonce":"1675065554644","channel":"envelope-channel","method":"invokeWithEnvelope","chaincode":"envelope-chaincode","deadline":"2033-01-31T07:58:39.677Z","public_key":"8JjYvYrzbeTuhuJnBJ7GKtwdofbnNwQnX8gmDrbNfYd2","signature":"36xWPPs7h1HKJgSGEHkzYqqP5M7gT44apSPya2RodBirPzsR7wrnvSXZu73rQnp4pJNYHKtVC3wBVcZkvfMmrnfk"}`
			b64Envelope := base64.StdEncoding.EncodeToString([]byte(jsonEnvelope))
			decodedEnvelope, err := e.DecodeEnvelope([]byte(b64Envelope))
			Expect(err).NotTo(HaveOccurred())

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(signer, chaincode)).WithChannel(channel)
			resp := envelopCC.Invoke(methodInvoke, payload, decodedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Disallow to verify signature with invalid payload", func() {
			serializedEnvelope, _ := createEnvelope(signer, payload, channel, chaincode, methodInvoke, deadline)

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(signer, chaincode)).WithChannel(channel)
			invalidPayload := []byte("invalid payload")
			resp := envelopCC.Invoke(methodInvoke, invalidPayload, serializedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 500))
		})

		It("Disallow to verify signature with invalid method", func() {
			serializedEnvelope, _ := createEnvelope(signer, payload, channel, chaincode, "invalid method", deadline)

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(signer, chaincode)).WithChannel(channel)
			resp := envelopCC.Invoke(methodInvoke, payload, serializedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 500))
		})

		It("Disallow to verify signature with invalid channel", func() {
			serializedEnvelope, _ := createEnvelope(signer, payload, "invalid channel", chaincode, methodInvoke, deadline)

			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(signer, chaincode)).WithChannel(channel)
			resp := envelopCC.Invoke(methodInvoke, payload, serializedEnvelope)

			Expect(resp.Status).To(BeNumerically("==", 500))
		})

		It("Don't check signature for query method", func() {
			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(signer, chaincode)).WithChannel(channel)
			resp := envelopCC.Query(methodQuery, payload)

			Expect(resp.Status).To(BeNumerically("==", 200))
		})

	})

	Describe("Nonce verification (replay attack)", func() {
		It("Disallow to execute tx with the same parameters (nonce, payload, pubkey)", func() {
			envelopCC = testcc.NewMockStub(chaincode, testdata.NewEnvelopCC(signer, chaincode)).WithChannel(channel)

			publicKey, privateKey, _ := ed25519.GenerateKey()
			nonce := "thesamenonce"
			hashToSign := signer.Hash(payload, nonce, channel, chaincode, methodInvoke, deadline.AsTime().Format(e.TimeLayout), publicKey)
			sig, _ := signer.Sign(payload, nonce, channel, chaincode, methodInvoke, deadline.AsTime().Format(e.TimeLayout), privateKey)
			envelope := &e.Envelope{
				PublicKey:  base58.Encode([]byte(publicKey)),
				Signature:  base58.Encode(sig),
				Nonce:      nonce,
				HashToSign: base58.Encode(hashToSign[:]),
				HashFunc:   "SHA256",
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

func createEnvelope(signer *e.Signer, payload []byte, channel, chaincode, method string, deadline ...*timestamppb.Timestamp) ([]byte, *e.Envelope) {
	publicKey, privateKey, _ := e.NewEd25519().GenerateKey()
	nonce := e.CreateNonce()

	envelope := &e.Envelope{
		PublicKey: base58.Encode([]byte(publicKey)),
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
	hashToSign := signer.Hash(payload, nonce, channel, chaincode, method, formatDeadline, publicKey)
	envelope.HashToSign = base58.Encode(hashToSign)

	sig, _ := signer.Sign(payload, nonce, channel, chaincode, method, formatDeadline, privateKey)
	envelope.Signature = base58.Encode(sig)

	serializedEnvelope, _ := serialize.PreferJSONSerializer.ToBytesFrom(envelope)
	return serializedEnvelope, envelope
}
