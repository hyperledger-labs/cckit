package envelope

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"

	"github.com/hyperledger-labs/cckit/extensions/envelope/crypto"
)

func Sign(crypto crypto.Crypto, payload []byte, nonce, channel, chaincode, method, deadline string, privateKey []byte) ([]byte, error) {
	pubKey, err := crypto.PublicKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf(`extract public key: %w`, err)
	}

	return crypto.Sign(privateKey,
		crypto.Hash(PrepareToHash(payload, nonce, channel, chaincode, method, deadline, pubKey)))
}

func CreateNonce() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func PrepareToHash(payload []byte, nonce, channel, chaincode, method, deadline string, pubkey []byte) []byte {
	bb := append(removeSpacesBetweenCommaAndQuotes(payload), nonce...) // resolve the unclear json serialization behavior in protojson package
	bb = append(bb, channel...)
	bb = append(bb, chaincode...)
	bb = append(bb, method...)
	bb = append(bb, deadline...)
	b58Pubkey := base58.Encode(pubkey)
	bb = append(bb, b58Pubkey...)
	return bb
}

func removeSpacesBetweenCommaAndQuotes(s []byte) []byte {
	removed := strings.ReplaceAll(string(s), `", "`, `","`)
	removed = strings.ReplaceAll(removed, `"}, {"`, `"},{"`)
	return []byte(strings.ReplaceAll(removed, `], "`, `],"`))
}
