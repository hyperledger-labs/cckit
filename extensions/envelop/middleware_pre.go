package envelop

import (
	"encoding/base64"
	"errors"
	"fmt"

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
				err := errors.New("signature not found")
				c.Logger().Sugar().Error(err)
				return router.ErrorResponse(err)
			} else {
				dst := make([]byte, base64.StdEncoding.DecodedLen(len(args[2])))
				n, err := base64.StdEncoding.Decode(dst, args[2])
				if err != nil {
					c.Logger().Error(`decod envelope failed:`, zap.Error(err))
					return router.ErrorResponse(err)
				}
				dst = dst[:n]
				fmt.Println("dst", string(dst))
				data, err := c.Serializer().FromBytesTo(dst, &Envelop{})
				fmt.Println("dst", string(dst))
				if err != nil {
					c.Logger().Error(`convert from bytes failed:`, zap.Error(err))
					return router.ErrorResponse(err)
				}
				envelope := data.(*Envelop)
				if err := CheckSig(args[1], envelope.Nonce, envelope.PublicKey, envelope.Signature); err != nil {
					c.Logger().Error(`check signature failed:`, zap.Error(err))
					return router.ErrorResponse(err)
				}
			}
		}
		return next(c)
	}
}
