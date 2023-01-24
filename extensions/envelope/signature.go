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

func Hash(payload []byte, nonce, channel, chaincode, method, deadline string, pubkey []byte) [32]byte {
	bb := append(payload, nonce...)
	bb = append(bb, channel...)
	bb = append(bb, chaincode...)
	bb = append(bb, method...)
	bb = append(bb, deadline...)
	bb = append(bb, pubkey...)
	return sha3.Sum256(bb)
}

func CreateSig(payload []byte, nonce, channel, chaincode, method, deadline string, privateKey []byte) ([]byte, []byte) {
	pubKey := ed25519.PrivateKey(privateKey).Public()
	hashed := Hash(payload, nonce, channel, chaincode, method, deadline, []byte(pubKey.(ed25519.PublicKey)))
	sig := ed25519.Sign(ed25519.PrivateKey(privateKey), hashed[:])
	return []byte(pubKey.(ed25519.PublicKey)), sig
}

func CheckSig(payload []byte, nonce, channel, chaincode, method, deadline string, pubKey []byte, sig []byte) error {
	hashed := Hash(payload, nonce, channel, chaincode, method, deadline, pubKey)
	if !ed25519.Verify(ed25519.PublicKey(pubKey), hashed[:], sig) {
		return ErrCheckSignatureFailed
	}
	return nil
}
