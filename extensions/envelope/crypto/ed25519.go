package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
)

func NewEd25519() *Ed25519 {
	return &Ed25519{}
}

type Ed25519 struct{}

func (ed *Ed25519) GenerateKey() (publicKey, privateKey []byte, err error) {
	publicKey, privateKey, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return publicKey, privateKey, nil
}

func (ed *Ed25519) Sign(privateKey, hash []byte) (signature []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("sign: %v", r)
		}
	}()
	return ed25519.Sign(privateKey, hash), nil
}

func (ed *Ed25519) Hash(msg []byte) []byte {
	h := sha256.Sum256(msg)
	return h[:]
}

func (ed *Ed25519) Verify(publicKey, hash, signature []byte) error {
	if !ed25519.Verify(publicKey, hash, signature) {
		return errors.New(`invalid signature`)
	}
	return nil
}

func (ed *Ed25519) PublicKey(privateKey []byte) ([]byte, error) {
	return ed25519.PrivateKey(privateKey).Public().(ed25519.PublicKey), nil
}
