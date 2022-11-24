package envelop

import "google.golang.org/protobuf/types/known/timestamppb"

type Envelop struct {
	PublicKey       []byte                 `json:"public_key,omitempty"`       // signer public key
	Signature       []byte                 `json:"signature,omitempty"`        // payload signature
	Nonce           string                 `json:"nonce,omitempty"`            // number is given for replay protection
	HashToSign      []byte                 `json:"hash_to_sign,omitempty"`     // payload hash
	HashFunc        string                 `json:"hash_func,omitempty"`        // function used for hashing
	Deadline        *timestamppb.Timestamp `json:"deadline,omitempty"`         // signature is not valid after deadline (EIP-2612)
	DomainSeparator []byte                 `json:"domain_separator,omitempty"` // to prevent replay attacks from other domains (EIP-2612)
}
