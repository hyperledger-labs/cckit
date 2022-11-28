package envelope

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hyperledger-labs/cckit/gateway"
	"google.golang.org/grpc/metadata"
)

// input opt for gateway to handle envelope with signature
func WithEnvelope() gateway.Opt {
	return func(opts *gateway.Opts) {
		opts.Input = append(opts.Input, func(ctx context.Context, input *gateway.ChaincodeInput) error {
			// get envelop with signature from header and add as second arg
			md, ok := metadata.FromIncomingContext(ctx)
			if ok {
				if v, ok := md["x-envelop"]; ok {
					envelope, err := DecodeEnvelope([]byte(v[0]))
					if err != nil {
						return fmt.Errorf(`invoke: %w`, err)
					}
					input.Args = append(input.Args, envelope)
				}
			}
			return nil
		})
	}
}

// decode base64 envelop
func DecodeEnvelope(encEnvelope []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(encEnvelope)))
	n, err := base64.StdEncoding.Decode(dst, encEnvelope)
	if err != nil {
		return nil, fmt.Errorf("parse envelope: %w", err)
	}
	return dst[:n], nil
}
