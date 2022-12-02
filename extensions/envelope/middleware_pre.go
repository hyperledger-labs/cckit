package envelope

import (
	"encoding/base64"
	"time"

	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger/fabric-protos-go/peer"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

const nonceObjectType = "nonce"

// pre-middleware for checking signature that is got in envelop
func Verify(next router.ContextHandlerFunc, pos ...int) router.ContextHandlerFunc {
	return func(c router.Context) peer.Response {
		iArgs := c.GetArgs()
		if len(iArgs) > 1 && iArgs[1] != nil {
			if len(iArgs) == 2 {
				c.Logger().Sugar().Error(ErrSignatureNotFound)
				return router.ErrorResponse(ErrSignatureNotFound)
			} else {
				data, err := c.Serializer().FromBytesTo(iArgs[2], &Envelope{})
				if err != nil {
					c.Logger().Error(`convert from bytes failed:`, zap.Error(err))
					return router.ErrorResponse(err)
				}
				envelope := data.(*Envelope)
				if envelope.Deadline.AsTime().Unix() < time.Now().Unix() {
					c.Logger().Sugar().Error(ErrDeadlineExpired)
					return router.ErrorResponse(ErrDeadlineExpired)
				}

				// check method and channel names because envelope can only be used once for channel+chaincode+method combination
				if string(iArgs[0]) != envelope.Method {
					c.Logger().Sugar().Error(ErrInvalidMethod)
					return router.ErrorResponse(ErrInvalidMethod)
				}
				if string(c.Stub().GetChannelID()) != envelope.Channel {
					c.Logger().Sugar().Error(ErrInvalidChannel)
					return router.ErrorResponse(ErrInvalidChannel)
				}

				// replay attack check
				txHash := txNonceKey(iArgs[1], envelope.Nonce, envelope.Channel, envelope.Chaincode, envelope.Method, envelope.PublicKey)
				key, err := c.Stub().CreateCompositeKey(nonceObjectType, []string{txHash})
				if err != nil {
					return router.ErrorResponse(err)
				}
				bb, err := c.Stub().GetState(key)
				if bb == nil && err == nil {
					if err := c.Stub().PutState(key, []byte{'0'}); err != nil {
						return router.ErrorResponse(err)
					}
					if err := CheckSig(iArgs[1], envelope.Nonce, envelope.Channel, envelope.Chaincode, envelope.Method, envelope.PublicKey, envelope.Signature); err != nil {
						c.Logger().Sugar().Error(ErrCheckSignatureFailed)
						return router.ErrorResponse(ErrCheckSignatureFailed)
					}
				} else {
					c.Logger().Sugar().Error(ErrTxAlreadyExecuted)
					return router.ErrorResponse(ErrTxAlreadyExecuted)
				}
			}
		}
		return next(c)
	}
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
