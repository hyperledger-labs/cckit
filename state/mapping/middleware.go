package mapping

import (
	"github.com/hyperledger-labs/cckit/router"
)

func MapStates(stateMappings StateMappings) router.MiddlewareFunc {
	return func(next router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(c router.Context) (interface{}, error) {
			c.UseState(WrapState(c.State(), stateMappings))
			return next(c)
		}
	}
}

func MapEvents(eventMappings EventMappings) router.MiddlewareFunc {
	return func(next router.HandlerFunc, pos ...int) router.HandlerFunc {
		return func(c router.Context) (interface{}, error) {
			c.UseEvent(WrapEvent(c.Event(), eventMappings))
			return next(c)
		}
	}
}
