package envelop

import (
	"time"

	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger/fabric-protos-go/peer"
	"go.uber.org/zap"
)

// pre-middleware for checking signature that is got in envelop
func Verify(next router.ContextHandlerFunc, pos ...int) router.ContextHandlerFunc {
	return func(c router.Context) peer.Response {
		args := c.GetArgs() // todo: check method type == invoke
		if len(args) > 1 {
			if len(args) == 2 {
				c.Logger().Sugar().Error(ErrSignatureNotFound)
				return router.ErrorResponse(ErrSignatureNotFound)
			} else {
				data, err := c.Serializer().FromBytesTo(args[2], &Envelop{})
				if err != nil {
					c.Logger().Error(`convert from bytes failed:`, zap.Error(err))
					return router.ErrorResponse(err)
				}
				env := data.(Envelop)
				if env.Deadline.AsTime().Unix() < time.Now().Unix() {
					c.Logger().Sugar().Error(ErrDeadlineExpired)
					return router.ErrorResponse(ErrDeadlineExpired)
				}
				if err := CheckSig(args[1], env.Nonce, env.PublicKey, env.Signature); err != nil {
					c.Logger().Sugar().Error(ErrCheckSignatureFailed)
					return router.ErrorResponse(ErrCheckSignatureFailed)
				}
			}
		}
		return next(c)
	}
}
