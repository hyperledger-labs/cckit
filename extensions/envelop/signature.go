package envelop

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

var (
	ErrCheckSignatureFailed = errors.New(`check signature failed`)
)

func CreatePrivKey() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return publicKey, privateKey, nil
}

func CreateNonce() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func Hash(payload []byte, nonce string) [32]byte {
	return sha3.Sum256(append(payload, nonce...))
}

func CreateSig(payload []byte, nonce string, privateKey []byte) ([]byte, []byte) {
	hashed := Hash(payload, nonce)
	pubKey := ed25519.PrivateKey(privateKey).Public()
	sig := ed25519.Sign(ed25519.PrivateKey(privateKey), hashed[:])
	return []byte(pubKey.(ed25519.PublicKey)), sig
}

func CheckSig(payload []byte, nonce string, pubKey []byte, sig []byte) error {
	hashed := Hash(payload, nonce)
	if !ed25519.Verify(ed25519.PublicKey(pubKey), hashed[:], sig) {
		return ErrCheckSignatureFailed
	}
	return nil
}
