package envelope

import (
	"encoding/base64"
	"time"

	"github.com/hyperledger-labs/cckit/router"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

const (
	// argument indexes
	methodNamePos = iota
	payloadPos
	sigPos

	nonceObjectType = "nonce"
	invokeType      = "invoke"
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
					if err := verifyEnvelope(c, iArgs[methodNamePos], iArgs[payloadPos], iArgs[sigPos]); err != nil {
						return nil, err
					}
				}
			}
			return next(c)
		}
	}
}

func verifyEnvelope(c router.Context, m, p, sig []byte) error {
	data, err := c.Serializer().FromBytesTo(sig, &Envelope{})
	if err != nil {
		c.Logger().Error(`convert from bytes failed:`, zap.Error(err))
		return err
	}
	envelope := data.(*Envelope)
	if envelope.Deadline.AsTime().Unix() < time.Now().Unix() {
		c.Logger().Sugar().Error(ErrDeadlineExpired)
		return ErrDeadlineExpired
	}

	// check method and channel names because envelope can only be used once for channel+chaincode+method combination
	if string(m) != envelope.Method {
		c.Logger().Sugar().Error(ErrInvalidMethod)
		return ErrInvalidMethod
	}
	if string(c.Stub().GetChannelID()) != envelope.Channel {
		c.Logger().Sugar().Error(ErrInvalidChannel)
		return ErrInvalidChannel
	}

	// replay attack check
	txHash := txNonceKey(p, envelope.Nonce, envelope.Channel, envelope.Chaincode, envelope.Method, envelope.PublicKey)
	key, err := c.Stub().CreateCompositeKey(nonceObjectType, []string{txHash})
	if err != nil {
		return err
	}
	bb, err := c.Stub().GetState(key)
	if bb == nil && err == nil {
		if err := c.Stub().PutState(key, []byte{'0'}); err != nil {
			return err
		}
		if err := CheckSig(p, envelope.Nonce, envelope.Channel, envelope.Chaincode, envelope.Method, envelope.Deadline.String(), envelope.PublicKey, envelope.Signature); err != nil {
			c.Logger().Sugar().Error(ErrCheckSignatureFailed)
			return ErrCheckSignatureFailed
		}
	} else {
		c.Logger().Sugar().Error(ErrTxAlreadyExecuted)
		return ErrTxAlreadyExecuted
	}
	return nil
}

func txNonceKey(payload []byte, nonce, channel, chaincode, method string, pubKey []byte) string {
	bb := append(payload, pubKey...)
	bb = append(bb, nonce...)
	bb = append(bb, channel...)
	bb = append(bb, chaincode...)
	bb = append(bb, method...)
	hashed := sha3.Sum256(bb)
	return base64.StdEncoding.EncodeToString(hashed[:])
}
