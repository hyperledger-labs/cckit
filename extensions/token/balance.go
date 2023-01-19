package token

import (
	"fmt"

	"github.com/hyperledger-labs/cckit/extensions/account"
	"github.com/hyperledger-labs/cckit/router"
)

type BalanceService struct {
	Account account.Getter
	Token   TokenGetter
	Store   BalanceStore
}

func NewBalanceService(
	accountResolver account.Getter,
	tokenGetter TokenGetter,
	store BalanceStore) *BalanceService {

	return &BalanceService{
		Account: accountResolver,
		Token:   tokenGetter,
		Store:   store,
	}
}

func (s *BalanceService) GetBalance(ctx router.Context, id *BalanceId) (*Balance, error) {
	if err := router.ValidateRequest(id); err != nil {
		return nil, err
	}

	_, err := s.Token.GetToken(ctx, &TokenId{Symbol: id.Symbol, Group: id.Group})
	if err != nil {
		return nil, fmt.Errorf(`get token: %w`, err)
	}
	return s.Store.Get(ctx, id)
}

func (s *BalanceService) ListBalances(ctx router.Context, id *BalanceId) (*Balances, error) {
	// empty balance id - no conditions
	balances, err := s.Store.List(ctx, &TokenId{
		Symbol: id.Symbol,
		Group:  id.Group,
	})
	if err != nil {
		return nil, err
	}

	return &Balances{Items: balances}, nil
}

func (s *BalanceService) Transfer(ctx router.Context, transfer *TransferRequest) (*TransferResponse, error) {
	if err := router.ValidateRequest(transfer); err != nil {
		return nil, err
	}

	invokerAddress, err := s.Account.GetInvokerAddress(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(`get invoker address: %w`, err)
	}

	_, err = s.Token.GetToken(ctx, &TokenId{Symbol: transfer.Symbol, Group: transfer.Group})
	if err != nil {
		return nil, fmt.Errorf(`get token: %w`, err)
	}

	if err := s.Store.Transfer(ctx, &TransferOperation{
		Sender:    invokerAddress.Address,
		Recipient: transfer.Recipient,
		Symbol:    transfer.Symbol,
		Group:     transfer.Group,
		Amount:    transfer.Amount,
		Meta:      transfer.Meta,
	}); err != nil {
		return nil, err
	}

	if err = Event(ctx).Set(&Transferred{
		Sender:    invokerAddress.Address,
		Recipient: transfer.Recipient,
		Symbol:    transfer.Symbol,
		Group:     transfer.Group,
		Amount:    transfer.Amount,
	}); err != nil {
		return nil, err
	}

	return &TransferResponse{
		Sender:    invokerAddress.Address,
		Recipient: transfer.Recipient,
		Symbol:    transfer.Symbol,
		Group:     transfer.Group,
		Amount:    transfer.Amount,
	}, nil
}

func (s *BalanceService) TransferBatch(ctx router.Context, transferBatch *TransferBatchRequest) (*TransferBatchResponse, error) {
	if err := router.ValidateRequest(transferBatch); err != nil {
		return nil, err
	}

	invokerAddress, err := s.Account.GetInvokerAddress(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(`get invoker address: %w`, err)
	}

	// check only first
	_, err = s.Token.GetToken(ctx, &TokenId{
		Symbol: transferBatch.Transfers[0].Symbol, Group: transferBatch.Transfers[0].Group})
	if err != nil {
		return nil, fmt.Errorf(`get token: %w`, err)
	}

	var (
		operations []*TransferOperation
		events     []*Transferred
		responses  []*TransferResponse
	)

	for _, t := range transferBatch.Transfers {
		operations = append(operations, &TransferOperation{
			Sender:    invokerAddress.Address,
			Recipient: t.Recipient,
			Symbol:    t.Symbol,
			Group:     t.Group,
			Amount:    t.Amount,
			Meta:      t.Meta,
		})

		events = append(events, &Transferred{
			Sender:    invokerAddress.Address,
			Recipient: t.Recipient,
			Symbol:    t.Symbol,
			Group:     t.Group,
			Amount:    t.Amount,
		})

		responses = append(responses, &TransferResponse{
			Sender:    invokerAddress.Address,
			Recipient: t.Recipient,
			Symbol:    t.Symbol,
			Group:     t.Group,
			Amount:    t.Amount,
		})
	}
	if err := s.Store.TransferBatch(ctx, operations); err != nil {
		return nil, err
	}

	if err = Event(ctx).Set(&TransferredBatch{
		Transfers: events,
	}); err != nil {
		return nil, err
	}

	return &TransferBatchResponse{
		Transfers: responses,
	}, nil
}
