package envelope

import (
	"github.com/hyperledger-labs/cckit/extensions/envelope/crypto"
)

type (
	Verifier interface {
		Verify(payload []byte, nonce, channel, chaincode, method, deadline string, publicKey []byte, sig []byte) error
	}

	DefaultVerifier struct {
		verifier crypto.Verifier
	}
)

func NewVerifier(verifier crypto.Verifier) *DefaultVerifier {
	return &DefaultVerifier{verifier: verifier}
}

func (s *DefaultVerifier) Verify(payload []byte, nonce, channel, chaincode, method, deadline string, pubKey []byte, sig []byte) error {
	if err := s.verifier.Verify(pubKey,
		s.verifier.Hash(PrepareToHash(payload, nonce, channel, chaincode, method, deadline, pubKey)), sig); err != nil {
		return ErrCheckSignatureFailed
	}
	return nil
}
