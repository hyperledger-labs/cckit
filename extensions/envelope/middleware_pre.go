package envelope

import (
	"time"

	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger/fabric-protos-go/peer"
	"go.uber.org/zap"
)

// pre-middleware for checking signature that is got in envelop
func Verify(next router.ContextHandlerFunc, pos ...int) router.ContextHandlerFunc {
	return func(c router.Context) peer.Response {
		iArgs := c.GetArgs()
		if len(iArgs) > 1 {
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
				if err := CheckSig(iArgs[1], envelope.Nonce, envelope.PublicKey, envelope.Signature); err != nil {
					c.Logger().Sugar().Error(ErrCheckSignatureFailed)
					return router.ErrorResponse(ErrCheckSignatureFailed)
				}
			}
		}
		return next(c)
	}
}
