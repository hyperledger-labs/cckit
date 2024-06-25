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
	"github.com/hyperledger-labs/cckit/extensions/envelope/crypto"
	"github.com/hyperledger-labs/cckit/serialize"
)

func TestEnvelop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Envelop suite")
}

var (
	deadline = timestamppb.New(time.Now().AddDate(0, 0, 2))
	payload  = []byte(`{"symbol":"GLD","decimals":"8","name":"Gold digital asset","type":"DM","underlying_asset":"gold","issuer_id":"GLDINC"}`)
)

var _ = Describe(`Envelope`, func() {

	c := crypto.NewEd25519()
	verifier := e.NewVerifier(c)

	Describe("Verifier", func() {
		It("Allow to verify valid signature", func() {
			nonce := e.CreateNonce()
			publicKey, privateKey, _ := c.GenerateKey()
			sig, err := e.Sign(c, payload, nonce, channel, chaincode, methodInvoke, deadline.String(), privateKey)
			Expect(err).NotTo(HaveOccurred())
			err = verifier.Verify(payload, nonce, channel, chaincode, methodInvoke, deadline.String(), publicKey, sig)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Disallow to verify signature with invalid payload", func() {
			nonce := e.CreateNonce()
			publicKey, privateKey, _ := c.GenerateKey()
			sig, _ := e.Sign(c, payload, nonce, channel, chaincode, methodInvoke, deadline.String(), privateKey)
			invalidPayload := []byte("invalid payload")
			err := verifier.Verify(invalidPayload, nonce, channel, chaincode, methodInvoke, deadline.String(), publicKey, sig)
			Expect(err).Should(MatchError(e.ErrSignatureCheckFailed))
		})
	})

	Describe("Handle base64 envelop", func() {

		It("Allow to parse base64 envelop", func() {
			_, envelope := createEnvelope(c, payload, channel, chaincode, methodInvoke, deadline)
			jj, _ := json.Marshal(envelope)
			b64 := base64.StdEncoding.EncodeToString(jj)
			bb, err := e.DecodeEnvelope([]byte(b64))
			Expect(err).NotTo(HaveOccurred())
			Expect(bb).To(Equal([]byte(jj)))
		})

	})

	Describe("Signature verification", func() {

		It("Allow to verify valid signature", func() {
			serializedEnvelope, _ := createEnvelope(c, payload, channel, chaincode, methodInvoke, deadline)
			resp := NewNewEnvelopCCMock(verifier).Invoke(methodInvoke, payload, serializedEnvelope)
			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Allow to verify valid signature without deadline", func() {
			serializedEnvelope, _ := createEnvelope(c, payload, channel, chaincode, methodInvoke)
			resp := NewNewEnvelopCCMock(verifier).Invoke(methodInvoke, payload, serializedEnvelope)
			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Allow to verify valid signature from the envelope in base64 format with deadline", func() {
			jsonEnvelope := `{"hash_func":"SHA256","hash_to_sign":"Bjd5D5qq3FFj1fSGt7QsHNqrfdDBzY7Cs5Wm45WZ19LE","nonce":"1675065805271","channel":"envelope-channel","method":"invokeWithEnvelope","chaincode":"envelope-chaincode","deadline":"2033-01-31T07:58:39.677Z","public_key":"EH9cLNxpg4FQmow9i1cF1b2vXkaJ17wUym7GSEeX6LQv","signature":"VG8P1ShwTGqVNhV8DVaZaVfxbZp7E8G9cpBbUPXFsGiwqLk1NJZFg2jt3ff1uJK92t2TbzkfyiL5ZTHsjYCdoQk"}`
			b64Envelope := base64.StdEncoding.EncodeToString([]byte(jsonEnvelope))
			decodedEnvelope, err := e.DecodeEnvelope([]byte(b64Envelope))
			Expect(err).NotTo(HaveOccurred())

			resp := NewNewEnvelopCCMock(verifier).Invoke(methodInvoke, payload, decodedEnvelope)
			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Allow to verify valid signature from the envelope in base64 format without deadline", func() {
			jsonEnvelope := `{"hash_func":"SHA256","hash_to_sign":"H4HmQKUQJm2bxJvSDPpHaP5vYYGvL5dUqWEfTzNN3eeH","nonce":"1675065554644","channel":"envelope-channel","method":"invokeWithEnvelope","chaincode":"envelope-chaincode","deadline":"2033-01-31T07:58:39.677Z","public_key":"8JjYvYrzbeTuhuJnBJ7GKtwdofbnNwQnX8gmDrbNfYd2","signature":"36xWPPs7h1HKJgSGEHkzYqqP5M7gT44apSPya2RodBirPzsR7wrnvSXZu73rQnp4pJNYHKtVC3wBVcZkvfMmrnfk"}`
			b64Envelope := base64.StdEncoding.EncodeToString([]byte(jsonEnvelope))
			decodedEnvelope, err := e.DecodeEnvelope([]byte(b64Envelope))
			Expect(err).NotTo(HaveOccurred())

			resp := NewNewEnvelopCCMock(verifier).Invoke(methodInvoke, payload, decodedEnvelope)
			Expect(resp.Status).To(BeNumerically("==", 200))
		})

		It("Disallow to verify signature with invalid payload", func() {
			serializedEnvelope, _ := createEnvelope(c, payload, channel, chaincode, methodInvoke, deadline)
			invalidPayload := []byte("invalid payload")

			resp := NewNewEnvelopCCMock(verifier).Invoke(methodInvoke, invalidPayload, serializedEnvelope)
			Expect(resp.Status).To(BeNumerically("==", 500))
		})

		It("Disallow to verify signature with invalid method", func() {
			serializedEnvelope, _ := createEnvelope(c, payload, channel, chaincode, "invalid method", deadline)

			resp := NewNewEnvelopCCMock(verifier).Invoke(methodInvoke, payload, serializedEnvelope)
			Expect(resp.Status).To(BeNumerically("==", 500))
		})

		It("Disallow to verify signature with invalid channel", func() {
			serializedEnvelope, _ := createEnvelope(c, payload, "invalid channel", chaincode, methodInvoke, deadline)

			resp := NewNewEnvelopCCMock(verifier).Invoke(methodInvoke, payload, serializedEnvelope)
			Expect(resp.Status).To(BeNumerically("==", 500))
		})

		It("Don't check signature for query method", func() {
			resp := NewNewEnvelopCCMock(verifier).Query(methodQuery, payload)
			Expect(resp.Status).To(BeNumerically("==", 200))
		})

	})

	Describe("Nonce verification (replay attack)", func() {
		It("Disallow to execute tx with the same parameters (nonce, payload, pubkey)", func() {
			publicKey, privateKey, _ := c.GenerateKey()
			nonce := "thesamenonce"

			hashToSign := c.Hash(e.PrepareToHash(payload, nonce, channel, chaincode, methodInvoke, deadline.AsTime().Format(e.TimeLayout), publicKey))
			sig, _ := c.Sign(privateKey, hashToSign)
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

			cc := NewNewEnvelopCCMock(verifier)
			resp := cc.Invoke(methodInvoke, payload, serializedEnvelope)
			Expect(resp.Status).To(BeNumerically("==", 200))

			resp = cc.Invoke(methodInvoke, payload, serializedEnvelope)
			Expect(errors.New(resp.Message)).To(MatchError(e.ErrTxAlreadyExecuted))
		})
	})

})

func createEnvelope(c crypto.Crypto, payload []byte, channel, chaincode, method string, deadline ...*timestamppb.Timestamp) ([]byte, *e.Envelope) {
	publicKey, privateKey, _ := crypto.NewEd25519().GenerateKey()
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
	hashToSign := c.Hash(e.PrepareToHash(payload, nonce, channel, chaincode, method, formatDeadline, publicKey))
	envelope.HashToSign = base58.Encode(hashToSign)

	sig, _ := c.Sign(privateKey, hashToSign)
	envelope.Signature = base58.Encode(sig)

	serializedEnvelope, _ := serialize.PreferJSONSerializer.ToBytesFrom(envelope)
	return serializedEnvelope, envelope
}
