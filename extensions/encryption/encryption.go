package encryption

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/serialize"
)

const TransientMapKey = `ENCODE_KEY`

// EncryptArgs convert args to [][]byte and encrypt args with key
func EncryptArgs(key []byte, args []interface{}, toBytesConverter serialize.ToBytesConverter) ([][]byte, error) {
	argBytes, err := serialize.ArgsToBytes(args, toBytesConverter)
	if err != nil {
		return nil, err
	}

	return EncryptArgsBytes(key, argBytes)
}

// EncryptArgsBytes encrypt args with key
func EncryptArgsBytes(key []byte, argsBytes [][]byte) ([][]byte, error) {
	eargs := make([][]byte, len(argsBytes))
	for i, bb := range argsBytes {
		encrypted, err := EncryptBytes(key, bb)
		if err != nil {
			return nil, errors.Wrap(err, `encryption error`)
		}

		eargs[i] = encrypted
	}
	return eargs, nil
}

// DecryptArgs decrypt args
func DecryptArgs(key []byte, args [][]byte) ([][]byte, error) {
	dargs := make([][]byte, len(args))
	for i, a := range args {

		// do not try to decrypt init function
		if i == 0 && string(a) == router.InitFunc {
			dargs[i] = a
			continue
		}

		decrypted, err := DecryptBytes(key, a)
		if err != nil {
			return nil, errors.Wrap(err, `decryption error`)
		}
		dargs[i] = decrypted
	}
	return dargs, nil
}

// Encrypt converts value to []byte  and encrypts its with key
func Encrypt(key []byte, value interface{}, toBytesConverter serialize.ToBytesConverter) ([]byte, error) {
	// TODO: customize  IV
	bb, err := toBytesConverter.ToBytesFrom(value)
	if err != nil {
		return nil, fmt.Errorf(`convert values to bytes: %w`, err)
	}
	return EncryptBytes(key, bb)
}

// Decrypt decrypts value with key
func Decrypt(key, value []byte) ([]byte, error) {
	return DecryptBytes(key, value)
}

func EncryptBytes(key, value []byte) ([]byte, error) {
	bb := make([]byte, len(value))
	copy(bb, value)
	// IV temporary blank
	ups, err := aesCBCEncryptWithIV(make([]byte, 16), key, pkcs7Padding(bb))
	return ups, err
}

func DecryptBytes(key, value []byte) ([]byte, error) {
	bb := make([]byte, len(value))
	copy(bb, value)
	return AESCBCPKCS7Decrypt(key, bb)
}

// TransientMapWithKey creates transient map with encrypting/decrypting key
func TransientMapWithKey(key []byte) map[string][]byte {
	return map[string][]byte{TransientMapKey: key}
}
