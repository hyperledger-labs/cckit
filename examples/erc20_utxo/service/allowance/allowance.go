package allowance

import (
	"errors"
	"fmt"

	"github.com/hyperledger-labs/cckit/extensions/account"
	"github.com/hyperledger-labs/cckit/extensions/token"
	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/state"
)

var (
	ErrOwnerOnly             = errors.New(`owner only`)
	ErrAllowanceInsufficient = errors.New(`allowance insufficient`)
)

type Service struct {
	balance token.BalanceStore
	account account.Getter
}

func NewService(account account.Getter, balance token.BalanceStore) *Service {
	return &Service{
		account: account,
		balance: balance,
	}
}

func (s *Service) GetAllowance(ctx router.Context, id *AllowanceId) (*Allowance, error) {
	if err := router.ValidateRequest(id); err != nil {
		return nil, err
	}

	allowance, err := State(ctx).Get(id, &Allowance{})
	if err != nil {
		if errors.Is(err, state.ErrKeyNotFound) {
			return &Allowance{
				Owner:   id.Owner,
				Spender: id.Spender,
				Symbol:  id.Symbol,
				Group:   id.Group,
				Amount:  nil,
			}, nil
		}
		return nil, fmt.Errorf(`get allowance: %w`, err)
	}

	return allowance.(*Allowance), nil
}

func (s *Service) Approve(ctx router.Context, approve *ApproveRequest) (*Allowance, error) {
	if err := router.ValidateRequest(approve); err != nil {
		return nil, err
	}

	invokerAddress, err := s.account.GetInvokerAddress(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(`get invoker address: %w`, err)
	}

	if invokerAddress.Address != approve.Owner {
		return nil, ErrOwnerOnly
	}

	allowance := &Allowance{
		Owner:   approve.Owner,
		Spender: approve.Spender,
		Symbol:  approve.Symbol,
		Group:   approve.Group,
		Amount:  approve.Amount,
	}

	if err := State(ctx).Put(allowance); err != nil {
		return nil, fmt.Errorf(`set allowance: %w`, err)
	}

	if err = Event(ctx).Set(&Approved{
		Owner:   approve.Owner,
		Spender: approve.Spender,
		Amount:  approve.Amount,
	}); err != nil {
		return nil, err
	}

	return allowance, nil
}

func (s *Service) TransferFrom(ctx router.Context, req *TransferFromRequest) (*TransferFromResponse, error) {
	if err := router.ValidateRequest(req); err != nil {
		return nil, err
	}

	spenderAddress, err := s.account.GetInvokerAddress(ctx, nil)
	if err != nil {
		return nil, err
	}

	allowance, err := s.GetAllowance(ctx, &AllowanceId{
		Owner:   req.Owner,
		Spender: spenderAddress.Address,
		Symbol:  req.Symbol,
		Group:   req.Group,
	})
	if err != nil {
		return nil, err
	}

	reqAmount, err := req.Amount.BigInt()
	if err != nil {
		return nil, fmt.Errorf(`req amount: %w`, err)
	}

	curAmount, err := allowance.Amount.BigInt()
	if err != nil {
		return nil, fmt.Errorf(`cur amount: %w`, err)
	}
	if curAmount.Cmp(reqAmount) == -1 {
		return nil, fmt.Errorf(`request trasfer amount=%s, allowance=%s: %w`,
			req.Amount, allowance.Amount, ErrAllowanceInsufficient)
	}

	allowance.Amount = token.BigIntSubAsDecimal(curAmount, reqAmount)

	// sub from allowance
	if err := State(ctx).Put(allowance); err != nil {
		return nil, fmt.Errorf(`update allowance: %w`, err)
	}

	if err = s.balance.Transfer(ctx, &token.TransferOperation{
		Sender:    req.Owner,
		Recipient: req.Recipient,
		Symbol:    req.Symbol,
		Group:     req.Group,
		Amount:    req.Amount,
	}); err != nil {
		return nil, err
	}

	if err = Event(ctx).Set(&TransferredFrom{
		Owner:     req.Owner,
		Spender:   spenderAddress.Address,
		Recipient: req.Recipient,
		Amount:    req.Amount,
	}); err != nil {
		return nil, err
	}

	return &TransferFromResponse{
		Owner:     req.Owner,
		Recipient: req.Recipient,
		Amount:    req.Amount,
	}, nil
}
