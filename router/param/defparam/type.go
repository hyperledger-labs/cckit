package defparam

import (
	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
)

func Proto(target interface{}, argPoss ...int) router.MiddlewareFunc {
	return param.Proto(router.DefaultParam, target, argPoss...)
}

func String(argPoss ...int) router.MiddlewareFunc {
	return param.String(router.DefaultParam, argPoss...)
}
