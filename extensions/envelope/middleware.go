package envelope

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/serialize"
	"go.uber.org/zap"
)

const (
	// argument indexes
	methodNamePos = iota
	payloadPos
	envelopePos

	nonceObjectType = "nonce"
	invokeType      = "invoke"

	TimeLayout = "2006-01-02T15:04:05.000Z"
)

// middleware for checking signature that is got in envelop
func Verify() router.MiddlewareFunc {
	return func(next router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(c router.Context) (interface{}, error) {
			if c.Handler().Type == invokeType {
				iArgs := c.GetArgs()
				if len(iArgs) == 2 {
					c.Logger().Sugar().Error(ErrSignatureNotFound)
					return nil, ErrSignatureNotFound
				} else {
					if err := verifyEnvelope(c, iArgs[methodNamePos], iArgs[payloadPos], iArgs[envelopePos]); err != nil {
						return nil, err
					}
				}
			}
			return next(c)
		}
	}
}

func verifyEnvelope(c router.Context, method, payload, envlp []byte) error {
	// parse json envelope format (json is original format for envelope from frontend)
	serializer := serialize.PreferJSONSerializer
	data, err := serializer.FromBytesTo(envlp, &Envelope{})
	if err != nil {
		c.Logger().Error(`convert from bytes failed:`, zap.Error(err))
		return err
	}
	envelope := data.(*Envelope)
	if envelope.Deadline.AsTime().Unix() != 0 {
		if envelope.Deadline.AsTime().Unix() < time.Now().Unix() {
			c.Logger().Sugar().Error(ErrDeadlineExpired)
			return ErrDeadlineExpired
		}
	}

	// check method and channel names because envelope can only be used once for channel+chaincode+method combination
	if string(method) != envelope.Method {
		c.Logger().Sugar().Error(ErrInvalidMethod)
		return ErrInvalidMethod
	}
	if string(c.Stub().GetChannelID()) != envelope.Channel {
		c.Logger().Sugar().Error(ErrInvalidChannel)
		return ErrInvalidChannel
	}

	// replay attack check
	txHash := txNonceKey(payload, envelope.Nonce, envelope.Channel, envelope.Chaincode, envelope.Method, envelope.PublicKey)
	key, err := c.Stub().CreateCompositeKey(nonceObjectType, []string{txHash})
	if err != nil {
		return err
	}
	bb, err := c.Stub().GetState(key)
	if bb == nil && err == nil {
		if err := c.Stub().PutState(key, []byte{'0'}); err != nil {
			return err
		}
		// convert public key and sig from hex string
		pubkey, err := hex.DecodeString(envelope.PublicKey)
		if err != nil {
			return err
		}
		sig, err := hex.DecodeString(envelope.Signature)
		if err != nil {
			return err
		}
		// convert deadline to frontend format
		var deadline string
		if envelope.Deadline != nil {
			deadline = envelope.Deadline.AsTime().Format(TimeLayout)
		}
		if err := CheckSig(payload, envelope.Nonce, envelope.Channel, envelope.Chaincode, envelope.Method, deadline, pubkey, sig); err != nil {
			c.Logger().Sugar().Error(ErrCheckSignatureFailed)
			return ErrCheckSignatureFailed
		}
	} else {
		c.Logger().Sugar().Error(ErrTxAlreadyExecuted)
		return ErrTxAlreadyExecuted
	}
	return nil
}

func txNonceKey(payload []byte, nonce, channel, chaincode, method, pubKey string) string {
	bb := append(payload, pubKey...)
	bb = append(bb, nonce...)
	bb = append(bb, channel...)
	bb = append(bb, chaincode...)
	bb = append(bb, method...)
	hashed := sha256.Sum256(bb)
	return base64.StdEncoding.EncodeToString(hashed[:])
}
