package token_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger-labs/cckit/extensions/token"
	"github.com/hyperledger-labs/cckit/serialize"
)

var (
	utxo1 = &token.UTXO{
		Symbol:  "a",
		Group:   "b",
		Address: "c",
		TxId:    "d",
		Amount: &token.BigInt{
			Data: `12345`,
		},
	}

	utxo1JSON               = []byte(`{"symbol":"a","group":"b","address":"c","tx_id":"d","amount":"12345"}`)
	utxo1JSON_WithBigIntObj = []byte(`{"symbol":"a","group":"b","address":"c","tx_id":"d","amount": { "data" : "12345" }}`)
)
var _ = Describe(`Bigint`, func() {

	It(`allow to serialize to json`, func() {

		json, err := serialize.PreferJSONSerializer.ToBytesFrom(utxo1)
		Expect(err).NotTo(HaveOccurred())

		Expect(json).To(Equal(utxo1JSON))
	})

	It(`allow to deserialize from json (short version)`, func() {
		obj, err := serialize.PreferJSONSerializer.FromBytesTo(utxo1JSON, &token.UTXO{})
		Expect(err).NotTo(HaveOccurred())

		Expect(obj.(*token.UTXO).String()).To(Equal(utxo1.String()))
	})

	It(`allow to deserialize from json (full version)`, func() {
		obj, err := serialize.PreferJSONSerializer.FromBytesTo(utxo1JSON_WithBigIntObj, &token.UTXO{})
		Expect(err).NotTo(HaveOccurred())

		Expect(obj.(*token.UTXO).String()).To(Equal(utxo1.String()))
	})
})
