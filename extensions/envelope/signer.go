package envelope

import (
	"crypto/ed25519"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

type Signer struct {
	crypto Crypto
}

func NewSigner(crypto Crypto) *Signer {
	return &Signer{crypto: crypto}
}

func removeSpacesBetweenCommaAndQuotes(s []byte) []byte {
	removed := strings.ReplaceAll(string(s), `", "`, `","`)
	removed = strings.ReplaceAll(removed, `"}, {"`, `"},{"`)
	return []byte(strings.ReplaceAll(removed, `], "`, `],"`))
}

func CreateNonce() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func (b *Signer) Hash(payload []byte, nonce, channel, chaincode, method, deadline string, pubkey []byte) []byte {
	bb := append(removeSpacesBetweenCommaAndQuotes(payload), nonce...) // resolve the unclear json serialization behavior in protojson package
	bb = append(bb, channel...)
	bb = append(bb, chaincode...)
	bb = append(bb, method...)
	bb = append(bb, deadline...)
	b58Pubkey := base58.Encode(pubkey)
	bb = append(bb, b58Pubkey...)
	return b.crypto.Hash(bb)
}

func (b *Signer) Sign(
	payload []byte, nonce, channel, chaincode, method, deadline string, privateKey []byte) ([]byte, error) {
	pubKey, err := b.crypto.PublicKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf(`extract public key: %w`, err)
	}

	hashed := b.Hash(payload, nonce, channel, chaincode, method, deadline, pubKey)
	return b.crypto.Sign(hashed, privateKey)
}

func (b *Signer) CheckSignature(payload []byte, nonce, channel, chaincode, method, deadline string, pubKey []byte, sig []byte) error {
	hashed := b.Hash(payload, nonce, channel, chaincode, method, deadline, pubKey)
	if !ed25519.Verify(pubKey, hashed[:], sig) {
		return ErrCheckSignatureFailed
	}
	return nil
}
