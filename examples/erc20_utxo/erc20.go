package erc20_utxo

import (
	"errors"
	"math/big"

	"github.com/hyperledger-labs/cckit/examples/erc20_utxo/service/allowance"
	"github.com/hyperledger-labs/cckit/examples/erc20_utxo/service/config"
	"github.com/hyperledger-labs/cckit/extensions/account"
	"github.com/hyperledger-labs/cckit/extensions/token"
	"github.com/hyperledger-labs/cckit/router"
)

var (
	Token = &token.CreateTokenTypeRequest{
		Name:        `SomeToken`,
		Symbol:      `@`,
		Decimals:    2,
		TotalSupply: token.NewDecimal(big.NewInt(10000000)),
	}
)

func NewChaincode() (*router.Chaincode, error) {
	r := router.New(`erc20`)

	// accountSvc resolves address as base58( invoker.Cert.PublicKey )
	accountSvc := account.NewLocalService()
	tokenSvc := token.NewTokenService()

	store := token.NewUTXOStore()
	// Balance management service
	balanceSvc := token.NewBalanceService(accountSvc, tokenSvc, store)
	// Allowance management service
	allowanceSvc := allowance.NewService(accountSvc, store)

	configSvc := config.NewService(tokenSvc)

	r.Init(func(ctx router.Context) (interface{}, error) {
		// add token definition to state if not exists
		t, err := token.CreateDefault(ctx, tokenSvc, Token)
		if err != nil {
			if errors.Is(err, token.ErrTokenAlreadyExists) {
				return nil, nil
			}
			return nil, err
		}

		// get chaincode instantiator address
		ownerAddress, err := accountSvc.GetInvokerAddress(ctx, nil)
		if err != nil {
			return nil, err
		}

		// add  `TotalSupply` to chaincode first committer
		if err = balanceSvc.Store.Mint(ctx, &token.BalanceOperation{
			Address: ownerAddress.Address,
			Symbol:  t.Symbol,
			Group:   t.Group,
			Amount:  Token.TotalSupply,
			Meta:    nil,
		}); err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err := token.RegisterBalanceServiceChaincode(r, balanceSvc); err != nil {
		return nil, err
	}
	if err := account.RegisterAccountServiceChaincode(r, accountSvc); err != nil {
		return nil, err
	}
	if err := config.RegisterConfigServiceChaincode(r, configSvc); err != nil {
		return nil, err
	}
	if err := allowance.RegisterAllowanceServiceChaincode(r, allowanceSvc); err != nil {
		return nil, err
	}

	return router.NewChaincode(r), nil
}
