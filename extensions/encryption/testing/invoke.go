package testing

import (
	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/hyperledger-labs/cckit/extensions/encryption"
	"github.com/hyperledger-labs/cckit/response"
	testcc "github.com/hyperledger-labs/cckit/testing"
)

// MockInvoke helper for invoking MockStub with transient key and encrypted args
func MockInvoke(cc *testcc.MockStub, encKey []byte, args ...interface{}) peer.Response {
	encArgs, err := encryption.EncryptArgs(encKey, args...)
	if err != nil {
		return response.Error(`unable to encrypt input args`)
	}
	return cc.AddTransient(encryption.TransientMapWithKey(encKey)).InvokeBytes(encArgs...)
}

// MockQuery helper for querying MockStub with transient key and encrypted args
func MockQuery(cc *testcc.MockStub, encKey []byte, args ...interface{}) peer.Response {
	encArgs, err := encryption.EncryptArgs(encKey, args...)
	if err != nil {
		return response.Error(`unable to encrypt input args`)
	}
	return cc.AddTransient(encryption.TransientMapWithKey(encKey)).QueryBytes(encArgs...)
}
