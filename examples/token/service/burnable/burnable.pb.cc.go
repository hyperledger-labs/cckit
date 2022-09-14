// Code generated by protoc-gen-cc-gateway. DO NOT EDIT.
// source: token/service/burnable/burnable.proto

/*
Package burnable contains
  *   chaincode methods names {service_name}Chaincode_{method_name}
  *   chaincode interface definition {service_name}Chaincode
  *   chaincode gateway definition {service_name}}Gateway
  *   chaincode service to cckit router registration func
*/
package burnable

import (
	context "context"
	_ "embed"
	errors "errors"

	cckit_gateway "github.com/hyperledger-labs/cckit/gateway"
	cckit_router "github.com/hyperledger-labs/cckit/router"
	cckit_defparam "github.com/hyperledger-labs/cckit/router/param/defparam"
	cckit_sdk "github.com/hyperledger-labs/cckit/sdk"
)

// BurnableServiceChaincode method names
const (

	// BurnableServiceChaincodeMethodPrefix allows to use multiple services with same method names in one chaincode
	BurnableServiceChaincodeMethodPrefix = ""

	BurnableServiceChaincode_Burn = BurnableServiceChaincodeMethodPrefix + "Burn"
)

// BurnableServiceChaincode chaincode methods interface
type BurnableServiceChaincode interface {
	Burn(cckit_router.Context, *BurnRequest) (*BurnResponse, error)
}

// RegisterBurnableServiceChaincode registers service methods as chaincode router handlers
func RegisterBurnableServiceChaincode(r *cckit_router.Group, cc BurnableServiceChaincode) error {

	r.Invoke(BurnableServiceChaincode_Burn,
		func(ctx cckit_router.Context) (interface{}, error) {
			return cc.Burn(ctx, ctx.Param().(*BurnRequest))
		},
		cckit_defparam.Proto(&BurnRequest{}))

	return nil
}

//go:embed burnable.swagger.json
var BurnableServiceSwagger []byte

// NewBurnableServiceGateway creates gateway to access chaincode method via chaincode service
func NewBurnableServiceGateway(sdk cckit_sdk.SDK, channel, chaincode string, opts ...cckit_gateway.Opt) *BurnableServiceGateway {
	return NewBurnableServiceGatewayFromInstance(
		cckit_gateway.NewChaincodeInstanceService(
			sdk,
			&cckit_gateway.ChaincodeLocator{Channel: channel, Chaincode: chaincode},
			opts...,
		))
}

func NewBurnableServiceGatewayFromInstance(chaincodeInstance cckit_gateway.ChaincodeInstance) *BurnableServiceGateway {
	return &BurnableServiceGateway{
		ChaincodeInstance: chaincodeInstance,
	}
}

// gateway implementation
// gateway can be used as kind of SDK, GRPC or REST server ( via grpc-gateway or clay )
type BurnableServiceGateway struct {
	ChaincodeInstance cckit_gateway.ChaincodeInstance
}

func (c *BurnableServiceGateway) Invoker() cckit_gateway.ChaincodeInstanceInvoker {
	return cckit_gateway.NewChaincodeInstanceServiceInvoker(c.ChaincodeInstance)
}

// ServiceDef returns service definition
func (c *BurnableServiceGateway) ServiceDef() cckit_gateway.ServiceDef {
	return cckit_gateway.NewServiceDef(
		_BurnableService_serviceDesc.ServiceName,
		BurnableServiceSwagger,
		&_BurnableService_serviceDesc,
		c,
		RegisterBurnableServiceHandlerFromEndpoint,
	)
}

func (c *BurnableServiceGateway) Burn(ctx context.Context, in *BurnRequest) (*BurnResponse, error) {
	var inMsg interface{} = in
	if v, ok := inMsg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}

	if res, err := c.Invoker().Invoke(ctx, BurnableServiceChaincode_Burn, []interface{}{in}, &BurnResponse{}); err != nil {
		return nil, err
	} else {
		return res.(*BurnResponse), nil
	}
}

// BurnableServiceChaincodeResolver interface for service resolver
type (
	BurnableServiceChaincodeResolver interface {
		Resolve(ctx cckit_router.Context) (BurnableServiceChaincode, error)
	}

	BurnableServiceChaincodeLocalResolver struct {
		service BurnableServiceChaincode
	}

	BurnableServiceChaincodeLocatorResolver struct {
		locatorResolver cckit_gateway.ChaincodeLocatorResolver
		service         BurnableServiceChaincode
	}
)

func NewBurnableServiceChaincodeLocalResolver(service BurnableServiceChaincode) *BurnableServiceChaincodeLocalResolver {
	return &BurnableServiceChaincodeLocalResolver{
		service: service,
	}
}

func (r *BurnableServiceChaincodeLocalResolver) Resolve(ctx cckit_router.Context) (BurnableServiceChaincode, error) {
	if r.service == nil {
		return nil, errors.New("service not set for local chaincode resolver")
	}

	return r.service, nil
}

func NewBurnableServiceChaincodeResolver(locatorResolver cckit_gateway.ChaincodeLocatorResolver) *BurnableServiceChaincodeLocatorResolver {
	return &BurnableServiceChaincodeLocatorResolver{
		locatorResolver: locatorResolver,
	}
}

func (r *BurnableServiceChaincodeLocatorResolver) Resolve(ctx cckit_router.Context) (BurnableServiceChaincode, error) {
	if r.service != nil {
		return r.service, nil
	}

	locator, err := r.locatorResolver(ctx, _BurnableService_serviceDesc.ServiceName)
	if err != nil {
		return nil, err
	}

	r.service = NewBurnableServiceChaincodeStubInvoker(locator)
	return r.service, nil
}

type BurnableServiceChaincodeStubInvoker struct {
	Invoker cckit_gateway.ChaincodeStubInvoker
}

func NewBurnableServiceChaincodeStubInvoker(locator *cckit_gateway.ChaincodeLocator) *BurnableServiceChaincodeStubInvoker {
	return &BurnableServiceChaincodeStubInvoker{
		Invoker: &cckit_gateway.LocatorChaincodeStubInvoker{Locator: locator},
	}
}

func (c *BurnableServiceChaincodeStubInvoker) Burn(ctx cckit_router.Context, in *BurnRequest) (*BurnResponse, error) {

	return nil, cckit_gateway.ErrInvokeMethodNotAllowed

}
