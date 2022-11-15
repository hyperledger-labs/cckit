// Package serialize for transforming  between json serialized  []byte and go structs
package serialize

const TypeInt = 1
const TypeString = ``
const TypeBool = true

type (
	// Serializer generic interface for serializing
	Serializer interface {
		ToBytesConverter
		FromBytesConverter
	}

	Serializable interface {
		ToBytes(ToBytesConverter) ([]byte, error)
	}

	// FromByter  interface supports FromBytes func for converting from slice of bytes to some defined target type
	FromByter interface {
		FromBytes([]byte) (interface{}, error)
	}

	// ToByter interface supports ToBytes func for converting to slice of bytes from source type
	ToByter interface {
		ToBytes() ([]byte, error)
	}

	// FromBytesConverter interface supports FromBytesConverter func for converting from slice of bytes to target type
	FromBytesConverter interface {
		FromBytesTo(from []byte, target interface{}) (interface{}, error)
	}

	// ToBytesConverter supports ToBytesConverter func converting from some interface to bytes
	ToBytesConverter interface {
		ToBytesFrom(from interface{}) ([]byte, error)
	}
)
