package envelope

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const keyEnvelope = "X-Envelope"

// handle custom header to pass envelope
func EnvelopeMatcher(key string) (string, bool) {
	switch key {
	case keyEnvelope:
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
