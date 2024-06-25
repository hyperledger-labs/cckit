package crypto

type (
	Crypto interface {
		Signer
		Hasher
		Verifier

		GenerateKey() (publicKey, privateKey []byte, err error)
		PublicKey(privateKey []byte) ([]byte, error)
	}

	Signer interface {
		Sign(privateKey, hash []byte) ([]byte, error)
	}

	Hasher interface {
		Hash([]byte) []byte
	}

	Verifier interface {
		Verify(publicKey, hash, signature []byte) error
		Hasher
	}
)
