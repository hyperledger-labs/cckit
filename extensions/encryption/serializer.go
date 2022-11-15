package encryption

import (
	"fmt"

	"github.com/hyperledger-labs/cckit/serialize"
)

type (
	Serializer struct {
		// serializer without encryption
		serializer serialize.Serializer
		key        []byte
	}
)

var _ serialize.Serializer = &Serializer{}

// FromBytesTo used for decrypting data after reading from state or receiving as argument
func (s *Serializer) FromBytesTo(from []byte, target interface{}) (interface{}, error) {
	decrypted, err := Decrypt(s.key, from)
	if err != nil {
		return nil, fmt.Errorf(`decrypt: %w`, err)
	}

	return s.serializer.FromBytesTo(decrypted, target)

}

func (s *Serializer) ToBytesFrom(from interface{}) ([]byte, error) {
	bb, err := s.serializer.ToBytesFrom(from)
	if err != nil {
		return nil, err
	}
	return EncryptBytes(s.key, bb)
}
