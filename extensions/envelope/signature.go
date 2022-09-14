package envelope

import (
	"crypto/ed25519"
	"errors"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

var (
	ErrCheckSignatureFailed = errors.New(`check signature failed`)
)

func CreateNonce() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func Hash(fn, payload, nonce string) [32]byte {
	return sha3.Sum256([]byte(fn + payload + nonce))
}

func CreateSignature(fn, payload string, nonce string, privateKey []byte) (pubKey, signature []byte) {
	//hashed := Hash(fn, payload, nonce)
	//
	//public := ed25519.PrivateKey(privateKey).Public()
	//
	//return nil, nil
	panic(`implement me`)
}

func VerifySignature(fn, payload, nonce string, pubKey []byte, signature []byte) error {

	hashed := Hash(fn, payload, nonce)

	if !ed25519.Verify(ed25519.PublicKey(pubKey), hashed[:], signature) {
		return ErrCheckSignatureFailed
	}

	return nil
}
