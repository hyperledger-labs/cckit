package envelope

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const HeaderKey = "X-Envelope"

// handle custom header to pass envelope
func EnvelopeMatcher(key string) (string, bool) {
	switch key {
	case HeaderKey:
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
