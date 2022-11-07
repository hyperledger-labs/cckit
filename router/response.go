package router

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"go.uber.org/zap"

	"github.com/hyperledger-labs/cckit/serialize"
)

type ContextResponse struct {
	Context Context
}

type Response interface {
	Error(err interface{}) peer.Response
	Success(data interface{}) peer.Response
	Create(data interface{}, err interface{}) peer.Response
}

// Error response
func (c *ContextResponse) Error(err interface{}) peer.Response {
	res := ErrorResponse(err)
	c.Context.Logger().Error(`router handler error`,
		zap.String(`path`, c.Context.Path()),
		zap.String(`message`, res.Message))
	return res
}

// Success response
func (c *ContextResponse) Success(data interface{}) peer.Response {
	res := SuccessResponse(data, c.Context.Serializer())
	c.Context.Logger().Debug(`route handle success`,
		zap.String(`path`, c.Context.Path()),
		zap.ByteString(`data`, res.Payload))
	return res
}

// Create  returns error response if err != nil
func (c *ContextResponse) Create(data interface{}, err interface{}) peer.Response {
	result := CreateResponse(data, err, c.Context.Serializer())

	if result.Status == shim.ERROR {
		return c.Error(result.Message)
	}
	return c.Success(result.Payload)
}

// ErrorResponse returns shim.Error
func ErrorResponse(err interface{}) peer.Response {
	return shim.Error(fmt.Sprintf("%s", err))
}

// SuccessResponse  returns shim.Success with serialized json if necessary
func SuccessResponse(data interface{}, toBytesConverter serialize.ToBytesConverter) peer.Response {
	bb, err := toBytesConverter.ToBytesFrom(data)
	if err != nil {
		return shim.Success(nil)
	}
	return shim.Success(bb)
}

// CreateResponse  returns peer.Response (Success or Error) depending on value of err
// if err is (bool) false or is error interface - returns shim.Error
func CreateResponse(data interface{}, err interface{}, toBytesConverter serialize.ToBytesConverter) peer.Response {
	var errObj error

	switch e := err.(type) {
	case nil:
		errObj = nil
	case bool:
		if !e {
			errObj = errors.New(`boolean error: false`)
		}
	case string:
		if e != `` {
			errObj = errors.New(e)
		}
	case error:
		errObj = e
	default:
		panic(fmt.Sprintf(`unknowm error type %s`, err))

	}

	if errObj != nil {
		return ErrorResponse(errObj)
	}
	return SuccessResponse(data, toBytesConverter)
}
