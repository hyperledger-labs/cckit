package envelope

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

type (
	Signer interface {
		CreateNonce() string
		Hash(payload []byte, nonce, channel, chaincode, method, deadline string, publicKey []byte) []byte
		Sign(payload []byte, nonce, channel, chaincode, method, deadline string, privateKey []byte) ([]byte, error)
		CheckSignature(payload []byte, nonce, channel, chaincode, method, deadline string, publicKey []byte, sig []byte) error
		Crypto() Crypto
	}

	DefaultSigner struct {
		crypto Crypto
	}
)

func NewSigner(crypto Crypto) *DefaultSigner {
	return &DefaultSigner{crypto: crypto}
}

func removeSpacesBetweenCommaAndQuotes(s []byte) []byte {
	removed := strings.ReplaceAll(string(s), `", "`, `","`)
	removed = strings.ReplaceAll(removed, `"}, {"`, `"},{"`)
	return []byte(strings.ReplaceAll(removed, `], "`, `],"`))
}

func (s *DefaultSigner) CreateNonce() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func (s *DefaultSigner) PrepareToHash(payload []byte, nonce, channel, chaincode, method, deadline string, pubkey []byte) []byte {
	bb := append(removeSpacesBetweenCommaAndQuotes(payload), nonce...) // resolve the unclear json serialization behavior in protojson package
	bb = append(bb, channel...)
	bb = append(bb, chaincode...)
	bb = append(bb, method...)
	bb = append(bb, deadline...)
	b58Pubkey := base58.Encode(pubkey)
	bb = append(bb, b58Pubkey...)
	return bb
}

func (s *DefaultSigner) Hash(payload []byte, nonce, channel, chaincode, method, deadline string, pubkey []byte) []byte {
	return s.crypto.Hash(s.PrepareToHash(payload, nonce, channel, chaincode, method, deadline, pubkey))
}

func (s *DefaultSigner) Sign(
	payload []byte, nonce, channel, chaincode, method, deadline string, privateKey []byte) ([]byte, error) {
	pubKey, err := s.crypto.PublicKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf(`extract public key: %w`, err)
	}

	hashed := s.Hash(payload, nonce, channel, chaincode, method, deadline, pubKey)
	return s.crypto.Sign(privateKey, hashed)
}

func (s *DefaultSigner) CheckSignature(payload []byte, nonce, channel, chaincode, method, deadline string, pubKey []byte, sig []byte) error {
	hashed := s.Hash(payload, nonce, channel, chaincode, method, deadline, pubKey)
	if err := s.crypto.Verify(pubKey, hashed, sig); err != nil {
		return ErrCheckSignatureFailed
	}
	return nil
}
func (s *DefaultSigner) Crypto() Crypto {
	return s.crypto
}
