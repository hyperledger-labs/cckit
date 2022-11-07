package expect

import (
	"github.com/hyperledger-labs/cckit/serialize"
)

func defaultFromBytesConverter(fromBytesConverters ...serialize.FromBytesConverter) serialize.FromBytesConverter {
	if len(fromBytesConverters) > 0 {
		return fromBytesConverters[0]
	}

	return serialize.DefaultSerializer
}
