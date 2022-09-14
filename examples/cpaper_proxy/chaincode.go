package cpaper_proxy

import (
	"github.com/hyperledger-labs/cckit/examples/cpaper_asservice"
	"github.com/hyperledger-labs/cckit/extensions/crosscc"
	"github.com/hyperledger-labs/cckit/router"
)

func NewCCWithLocalCpaper() (*router.Chaincode, error) {
	r := router.New(`crosscc_local`)

	cpaperService := cpaper_asservice.NewService()
	crossCCService := NewServiceWithLocalCPaperResolver(cpaperService)

	// 2 services in one chaincode
	// both CPaper and CrossCC in one chaincode, that is why used local resolver
	if err := cpaper_asservice.RegisterCPaperServiceChaincode(r, cpaperService); err != nil {
		return nil, err
	}

	if err := RegisterCPaperProxyServiceChaincode(r, crossCCService); err != nil {
		return nil, err
	}

	return router.NewChaincode(r), nil
}

func NewCCWithRemoteCpaper() (*router.Chaincode, error) {
	r := router.New(`crosscc_remote`)

	crossCCSettingService := crosscc.NewSettingService()
	crossCCService := NewServiceWithRemoteCPaperResolver(crossCCSettingService)

	// crossCC service and CPaper service - in separate chaincodes
	// in crossCC chauincode there are two services:
	// 1. CrossCC itself
	// 2. Setting service to store information where (channel, chaincode) CPaper service located
	if err := crosscc.RegisterSettingServiceChaincode(r, crossCCSettingService); err != nil {
		return nil, err
	}

	if err := RegisterCPaperProxyServiceChaincode(r, crossCCService); err != nil {
		return nil, err
	}

	return router.NewChaincode(r), nil
}
