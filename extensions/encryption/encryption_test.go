package encryption_test

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/hyperledger/fabric-protos-go/peer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger-labs/cckit/examples/payment"
	"github.com/hyperledger-labs/cckit/examples/payment/schema"
	"github.com/hyperledger-labs/cckit/extensions/encryption"
	"github.com/hyperledger-labs/cckit/extensions/encryption/testdata"
	enctest "github.com/hyperledger-labs/cckit/extensions/encryption/testing"
	identitytestdata "github.com/hyperledger-labs/cckit/identity/testdata"
	"github.com/hyperledger-labs/cckit/serialize"
	"github.com/hyperledger-labs/cckit/state"
	"github.com/hyperledger-labs/cckit/state/mapping"
	testcc "github.com/hyperledger-labs/cckit/testing"
	expectcc "github.com/hyperledger-labs/cckit/testing/expect"
	"github.com/hyperledger-labs/cckit/testing/gomega"
)

func TestEncryption(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Router suite")
}

var (
	Owner = identitytestdata.Certificates[0].MustIdentity(`SOME_MSP`)
)

var (
	encryptOnDemandPaymentCC            *testcc.MockStub
	encryptPaymentCC                    *testcc.MockStub
	encryptPaymentCCWithEncStateContext *testcc.MockStub
	externalCC                          *testcc.MockStub
	serializer                          = serialize.DefaultSerializer

	encCCInvoker *enctest.MockStub

	pType  = `SALE`
	encKey []byte

	// fixtures
	pID1           = `id-1`
	pAmount1 int32 = 111

	pID2           = `id-2`
	pAmount2 int32 = 222

	pID3           = `id-3`
	pAmount3 int32 = 333

	encryptedPType, encryptedPId1   []byte
	payment1                        *schema.Payment
	payment1Serialized, encPayment1 []byte
	paymentMapper                   mapping.StateMapper
	encPaymentNamespace             []byte
	err                             error
)

var _ = Describe(`Encryption`, func() {

	BeforeSuite(func() {
		// Create encode key. In real case it can be calculated via ECDH
		encKey = make([]byte, 32)
		_, err = rand.Read(encKey)

		Expect(err).NotTo(HaveOccurred())

		//Create chaincode mock
		encryptOnDemandPaymentCC = testcc.NewMockStub(
			`paymentsEncOnDemand`,
			payment.NewEncryptOnDemandPaymentCC())

		encryptPaymentCC = testcc.NewMockStub(
			`paymentsEnc`,
			payment.NewEncryptPaymentCC())

		encryptPaymentCCWithEncStateContext = testcc.NewMockStub(
			`paymentsEncWithContext`,
			payment.NewEncryptedPaymentCCWithEncStateContext())

		encCCInvoker = enctest.NewMockStub(encryptPaymentCCWithEncStateContext, encKey)
		encCCInvoker.DecryptInvokeResponse = true

		externalCC = testcc.NewMockStub(`external`,
			testdata.NewExternalCC(`paymentsEncWithContext`, `payment-channel`))

		// external cc have access to encrypted payment chaincode
		externalCC.MockPeerChaincode(
			`paymentsEncWithContext/payment-channel`,
			encryptPaymentCCWithEncStateContext)

		encryptedPType, err = encryption.Encrypt(encKey, pType, encryptPaymentCC.Serializer)
		Expect(err).To(BeNil())

		encryptedPId1, err = encryption.Encrypt(encKey, pID1, encryptPaymentCC.Serializer)
		Expect(err).To(BeNil())

		payment1 = &schema.Payment{
			Type:   pType,
			Id:     pID1,
			Amount: pAmount1,
		}

		payment1Serialized, err = serializer.ToBytesFrom(payment1)
		Expect(err).To(BeNil())

		encPayment1, err = encryption.Encrypt(encKey, payment1, encryptPaymentCC.Serializer)
		Expect(err).To(BeNil())

		paymentMapper, _ = payment.StateMappings.Get(&schema.Payment{})
		Expect(err).To(BeNil())

		// we know that Payment namespace contain only one part
		encPaymentNamespace, err = encryption.Encrypt(encKey, paymentMapper.Namespace()[0], encryptPaymentCC.Serializer)
		Expect(err).To(BeNil())
	})

	Describe("Encrypting methods", func() {

		var encrypted1, encrypted2 []byte
		It("Allow to encrypt bytes", func() {
			encrypted1, err = encryption.EncryptBytes(encKey, []byte(pID1))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(encrypted1)).To(BeNumerically(">", len(pID1)))

			// once more
			encrypted2, err = encryption.EncryptBytes(encKey, []byte(pID1))
			Expect(err).NotTo(HaveOccurred())
			Expect(encrypted2).To(Equal(encrypted1))
		})

		It("Allow to decrypt bytes", func() {
			decrypted, err := encryption.DecryptBytes(encKey, encrypted1)
			Expect(err).NotTo(HaveOccurred())
			Expect(decrypted).To(Equal([]byte(pID1)))
		})

		It("Allow to encrypt bytes with serializer", func() {
			encrypted2, err = encryption.Encrypt(encKey, []byte(pID1), serializer)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(encrypted2)).To(BeNumerically(">", len(pID1)))
			Expect(encrypted2).To(Equal(encrypted1))
		})

		It("Allow to decrypt bytes with serializer", func() {
			decrypted, err := encryption.Decrypt(encKey, encrypted2)
			Expect(err).NotTo(HaveOccurred())
			Expect(decrypted).To(Equal([]byte(pID1)))
		})
		It("Allow to encrypt struct", func() {
			encrypted1, err = encryption.EncryptBytes(encKey, []byte(pID1))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(encrypted1)).To(BeNumerically(">", len(pID1)))
		})

		It("Allow to decrypt bytes", func() {
			decrypted, err := encryption.DecryptBytes(encKey, encrypted1)
			Expect(err).NotTo(HaveOccurred())
			Expect(decrypted).To(Equal([]byte(pID1)))
		})

		It(`Allow to encrypt args`, func() {

			args, err := encryption.EncryptArgs(encKey,
				[]interface{}{payment.PaymentGet, pType, pID1}, serializer)
			Expect(err).To(BeNil())
			Expect(len(args)).To(Equal(3))

			methodNameEncrypted, err := encryption.Encrypt(encKey, payment.PaymentGet, serializer)
			Expect(err).To(BeNil())
			Expect(args[0]).To(Equal(methodNameEncrypted))

			Expect(args[1]).To(Equal(encryptedPType))
		})

	})
	Describe("Encrypting in demand with state mapping", func() {

		It("Allow to init encrypt-on-demand payment chaincode without key in transient map", func() {
			expectcc.ResponseOk(encryptOnDemandPaymentCC.Init())
		})

		It("Disallow to create encrypted payment providing unencrypted arguments", func() {
			expectcc.ResponseError(
				encryptOnDemandPaymentCC.WithTransient(encryption.TransientMapWithKey(encKey)).
					Invoke(`paymentCreate`, pType, pID1, pAmount1), `args: decryption error`)
		})

		It("Allow to create encrypted payment", func() {
			// encrypt all arguments
			args, err := encryption.EncryptArgs(encKey,
				[]interface{}{`paymentCreate`, pType, pID1, pAmount1}, encryptOnDemandPaymentCC.Serializer)
			Expect(err).To(BeNil())
			Expect(len(args)).To(Equal(4))

			Expect(err).To(BeNil())
			Expect(args[1]).To(Equal(encryptedPType))
			// second argument is encoded payment id
			Expect(args[2]).To(Equal(encryptedPId1))
			// invoke chaincode with encoded args and encKey via transientMap, receives encoded payment id
			ccPId := expectcc.PayloadBytes(
				encryptOnDemandPaymentCC.WithTransient(encryption.TransientMapWithKey(encKey)).
					InvokeBytes(args...), encryptedPId1)

			encEvent := <-encryptOnDemandPaymentCC.ChaincodeEventsChannel
			decryptedEvent := encryption.MustDecryptEvent(encKey, encEvent)

			Expect(decryptedEvent.Payload).To(BeEquivalentTo(
				testcc.MustProtoMarshal(&schema.PaymentEvent{Type: pType, Id: pID1, Amount: pAmount1})))

			Expect(decryptedEvent.EventName).To(Equal(`PaymentEvent`))

			decryptedPaymentID, err := encryption.Decrypt(encKey, ccPId)
			Expect(err).To(BeNil())
			Expect(string(decryptedPaymentID)).To(Equal(pID1))
		})

		It("Allow to get encrypted payment by type and id", func() {
			//encrypt all arguments
			args, err := encryption.EncryptArgs(encKey,
				[]interface{}{`paymentGet`, pType, pID1}, encryptPaymentCC.Serializer)
			Expect(err).To(BeNil())

			// Check that value is encrypted in chaincode state - use debugStateGet func
			// without providing key in transient map - so we need to provide encrypted key
			// and cause we dont't require key - state also returns unencrypted
			stateGetResult := encryptOnDemandPaymentCC.Invoke(`debugStateGet`, []string{
				base64.StdEncoding.EncodeToString(encPaymentNamespace),
				base64.StdEncoding.EncodeToString(encryptedPType),
				base64.StdEncoding.EncodeToString(encryptedPId1)})

			expectcc.PayloadBytes(stateGetResult, encPayment1)

			//returns unencrypted
			paymentFromCC := expectcc.PayloadIs(
				encryptOnDemandPaymentCC.WithTransient(encryption.TransientMapWithKey(encKey)).InvokeBytes(args...),
				&schema.Payment{}).(*schema.Payment)

			Expect(paymentFromCC).To(gomega.StringerEqual(payment1))
		})

		It("Allow to get encrypted payments by type as unencrypted values", func() {
			// MockInvoke sets key in transient map and encrypt input arguments
			payments := expectcc.PayloadIs(
				enctest.MockInvoke(encryptOnDemandPaymentCC, encKey, `paymentList`, pType),
				&schema.PaymentList{}).(*schema.PaymentList)

			Expect(payments.Items).To(HaveLen(1))
			// Returned value is not encrypted
			Expect(payments.Items[0].Id).To(Equal(pID1))
		})

		It("Allow to invoke with non encrypted data", func() {
			expectcc.ResponseOk(encryptOnDemandPaymentCC.Invoke(`paymentCreate`, pType, pID2, pAmount2))
		})

		It("Allow to get non encrypted payment by type and id", func() {
			//returns unencrypted
			paymentFromCC := expectcc.PayloadIs(
				encryptOnDemandPaymentCC.Query(`paymentGet`, pType, pID2), &schema.Payment{}).(*schema.Payment)

			Expect(paymentFromCC.Id).To(Equal(pID2))
		})
	})

	Describe("Encrypting required", func() {

		It("Disallow to init without key in transient map", func() {
			expectcc.ResponseError(encryptPaymentCC.Init())
		})

		It("Allow to init with key in transient map", func() {
			expectcc.ResponseOk(encryptPaymentCC.WithTransient(encryption.TransientMapWithKey(encKey)).Init())
		})

		It("Disallow to create payment without providing key in encryptPaymentCC ", func() {
			expectcc.ResponseError(encryptPaymentCC.Invoke(`paymentCreate`, pType, pID3, pAmount3),
				encryption.ErrKeyNotDefinedInTransientMap)
		})

		It("Allow to create payment providing key in encryptPaymentCC ", func() {
			// encode all arguments
			expectcc.ResponseOk(enctest.MockInvoke(encryptPaymentCC, encKey, `paymentCreate`, pType, pID3, pAmount3))
		})
	})

	Describe("Encrypted state context", func() {
		It("Allow to init encrypt payment chaincodes", func() {
			expectcc.ResponseOk(encryptPaymentCCWithEncStateContext.WithTransient(encryption.TransientMapWithKey(encKey)).Init())
		})
		//
		It("Allow to create payment providing key in encryptPaymentCC ", func(done Done) {
			events, closer := encryptPaymentCCWithEncStateContext.EventSubscription()

			responsePayment := expectcc.PayloadIs(
				// encCCInvoker encrypts args before passing to cc invoke and pass key in transient map
				encCCInvoker.From(Owner).Invoke(`paymentCreate`, pType, pID1, pAmount3),
				&schema.Payment{}).(*schema.Payment)

			// we use encryption.MockStub DecryptInvokeResponse feature
			Expect(responsePayment.Id).To(Equal(pID1))
			Expect(responsePayment.Type).To(Equal(pType))
			Expect(responsePayment.Amount).To(Equal(pAmount3))

			//event name and payload is encrypted with key
			Expect(<-events).To(BeEquivalentTo(encryption.MustEncryptEvent(encKey, &peer.ChaincodeEvent{
				ChaincodeId: encCCInvoker.MockStub.Name,
				TxId:        encCCInvoker.MockStub.LastTxID,
				EventName:   `PaymentEvent`,
				Payload: testcc.MustProtoMarshal(&schema.PaymentEvent{
					Type:   pType,
					Id:     pID1,
					Amount: pAmount3,
				}),
			})))

			_ = closer()
			close(done)
		}, 0.2)

		It("Allow to get payment by type and id", func() {
			// encCCInvoker encrypts args before passing to cc invoke and pass key in transient map
			paymentFromCC := expectcc.PayloadIs(
				encCCInvoker.From(Owner).Query(`paymentGet`, pType, pID1), &schema.Payment{}).(*schema.Payment)

			//returned payload is unencrypted
			Expect(paymentFromCC.Id).To(Equal(pID1))
			Expect(paymentFromCC.Type).To(Equal(pType))
			Expect(paymentFromCC.Amount).To(Equal(pAmount3))
		})

		It("Allow to get payment providing key using debugStateGet", func() {
			// we didn't provide encrypting key,
			// chaincode use ArgsDecrypt middleware, requiring key in transient map
			// but we add exception for method debugStateGet
			// for all chaincode methods, except stateGet @see example/payments/cc_enc_context_with_mapping.go
			expectcc.ResponseOk(encCCInvoker.MockStub.Invoke(`debugStateGet`, []string{
				base64.StdEncoding.EncodeToString(encPaymentNamespace),
				base64.StdEncoding.EncodeToString(encryptedPType),
				base64.StdEncoding.EncodeToString(encryptedPId1)}))
		})

		It("Allow to get encrypted payments by type as unencrypted values", func() {
			paymentsFromCC := expectcc.PayloadIs(
				encCCInvoker.Query(`paymentList`, pType), &schema.PaymentList{}).(*schema.PaymentList)

			Expect(paymentsFromCC.Items).To(HaveLen(1))
			// Returned value is not encrypted
			Expect(paymentsFromCC.Items[0].Id).To(Equal(pID1))
			Expect(paymentsFromCC.Items[0].Type).To(Equal(pType))
			Expect(paymentsFromCC.Items[0].Amount).To(Equal(pAmount3))
		})

		It("Disallow to get payment by type and id without providing encrypting key in transient map", func() {
			expectcc.ResponseError(encCCInvoker.MockStub.From(Owner).Query(`paymentGet`, pType, pID1),
				encryption.ErrKeyNotDefinedInTransientMap)
		})

		It("Disallow to get non existent payment by type and id providing encrypting key in transient map", func() {
			// key in error is not encrypted
			expectcc.ResponseError(encCCInvoker.From(Owner).Query(`paymentGet`, pType, pID1+`NoExists`),
				`Payment | SALE | id-1NoExists: `+state.ErrKeyNotFound.Error())
		})

		It("Allow to get payment by type and id", func() {
			//returns unencrypted
			paymentFromCC := expectcc.PayloadIs(encCCInvoker.From(Owner).Query(`paymentGet`, pType, pID1),
				&schema.Payment{}).(*schema.Payment)
			Expect(paymentFromCC.Id).To(Equal(pID1))
			Expect(paymentFromCC.Type).To(Equal(pType))
			Expect(paymentFromCC.Amount).To(Equal(pAmount3))
		})

		It("Allow to get payment via external chaincode", func() {
			paymentFromExtCC := expectcc.PayloadIs(externalCC.WithTransient(encryption.
				TransientMapWithKey(encKey)).Query(`checkPayment`, pType, pID1), &schema.Payment{}).(*schema.Payment)
			Expect(paymentFromExtCC.Id).To(Equal(pID1))
			Expect(paymentFromExtCC.Type).To(Equal(pType))
			Expect(paymentFromExtCC.Amount).To(Equal(pAmount3))
		})

	})
})
