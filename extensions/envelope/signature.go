package envelope

import (
	"crypto/ed25519"
	"crypto/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

func CreateKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return publicKey, privateKey, nil
}

func CreateNonce() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func Hash(payload []byte, nonce, channel, chaincode, method string) [32]byte {
	bb := append(payload, nonce...)
	bb = append(bb, channel...)
	bb = append(bb, chaincode...)
	bb = append(bb, method...)
	return sha3.Sum256(bb)
}

func CreateSig(payload []byte, nonce, channel, chaincode, method string, privateKey []byte) ([]byte, []byte) {
	hashed := Hash(payload, nonce, channel, chaincode, method)
	pubKey := ed25519.PrivateKey(privateKey).Public()
	sig := ed25519.Sign(ed25519.PrivateKey(privateKey), hashed[:])
	return []byte(pubKey.(ed25519.PublicKey)), sig
}

func CheckSig(payload []byte, nonce, channel, chaincode, method string, pubKey []byte, sig []byte) error {
	hashed := Hash(payload, nonce, channel, chaincode, method)
	if !ed25519.Verify(ed25519.PublicKey(pubKey), hashed[:], sig) {
		return ErrCheckSignatureFailed
	}
	return nil
}
