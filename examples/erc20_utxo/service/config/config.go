package config

import (
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hyperledger-labs/cckit/extensions/token"
	"github.com/hyperledger-labs/cckit/router"
)

type Service struct {
	Token token.TokenGetter
}

func NewService(tokenGetter token.TokenGetter) *Service {
	return &Service{
		Token: tokenGetter,
	}
}

func (s *Service) GetName(ctx router.Context, e *emptypb.Empty) (*NameResponse, error) {
	t, err := s.Token.GetDefaultToken(ctx, e)
	if err != nil {
		return nil, err
	}

	return &NameResponse{Name: t.GetType().GetName()}, nil
}

func (s *Service) GetSymbol(ctx router.Context, e *emptypb.Empty) (*SymbolResponse, error) {
	t, err := s.Token.GetDefaultToken(ctx, e)
	if err != nil {
		return nil, err
	}

	return &SymbolResponse{Symbol: t.GetType().GetSymbol()}, nil
}

func (s *Service) GetDecimals(ctx router.Context, e *emptypb.Empty) (*DecimalsResponse, error) {
	t, err := s.Token.GetDefaultToken(ctx, e)
	if err != nil {
		return nil, err
	}

	return &DecimalsResponse{Decimals: t.GetType().GetDecimals()}, nil
}

func (s *Service) GetTotalSupply(ctx router.Context, e *emptypb.Empty) (*TotalSupplyResponse, error) {
	t, err := s.Token.GetDefaultToken(ctx, e)
	if err != nil {
		return nil, err
	}

	return &TotalSupplyResponse{TotalSupply: t.GetType().GetTotalSupply()}, nil
}
