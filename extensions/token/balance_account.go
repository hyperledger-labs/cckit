package token

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/hyperledger-labs/cckit/extensions/token/decimal"
	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/state"
)

var _ BalanceStore = &AccountStore{}

type AccountStore struct {
}

func NewAccountStore() *AccountStore {
	return &AccountStore{}
}

func (s *AccountStore) Get(ctx router.Context, id *BalanceId) (*Balance, error) {
	if err := router.ValidateRequest(id); err != nil {
		return nil, err
	}
	balance, err := State(ctx).Get(id, &Balance{})
	if err != nil {
		if strings.Contains(err.Error(), state.ErrKeyNotFound.Error()) {
			// default zero balance even if no Balance state entry for auth exists
			return &Balance{
				Address: id.Address,
				Symbol:  id.Symbol,
				Group:   id.Group,
				Amount:  decimal.New(big.NewInt(0)),
			}, nil
		}
		return nil, err
	}

	return balance.(*Balance), nil
}

func (s *AccountStore) List(ctx router.Context, id *TokenId) ([]*Balance, error) {
	balances, err := State(ctx).List(&Balance{})
	if err != nil {
		return nil, err
	}
	return balances.(*Balances).Items, nil
}

func (s *AccountStore) Transfer(ctx router.Context, transfer *TransferOperation) error {
	if err := router.ValidateRequest(transfer); err != nil {
		return err
	}
	// subtract from sender balance
	if _, err := s.sub(ctx, &BalanceOperation{
		Address: transfer.Sender,
		Symbol:  transfer.Symbol,
		Group:   transfer.Group,
		Amount:  transfer.Amount,
	}); err != nil {
		return fmt.Errorf(`subtract from sender: %w`, err)
	}
	// add to recipient balance
	if _, err := s.add(ctx, &BalanceOperation{
		Address: transfer.Recipient,
		Symbol:  transfer.Symbol,
		Group:   transfer.Group,
		Amount:  transfer.Amount,
	}); err != nil {
		return fmt.Errorf(`add to recipient: %w`, err)
	}
	return nil
}

func (s *AccountStore) TransferBatch(ctx router.Context, transfers []*TransferOperation) error {
	// todo: COUNT TOTAL AMOUNT !!!
	for _, t := range transfers {
		if err := s.Transfer(ctx, t); err != nil {
			return err
		}
	}

	return nil
}

func (s *AccountStore) Mint(ctx router.Context, op *BalanceOperation) error {
	if err := router.ValidateRequest(op); err != nil {
		return err
	}
	// add to recipient balance
	if _, err := s.add(ctx, op); err != nil {
		return fmt.Errorf(`add: %w`, err)
	}

	return nil
}

func (s *AccountStore) Burn(ctx router.Context, op *BalanceOperation) error {
	if err := router.ValidateRequest(op); err != nil {
		return err
	}
	if _, err := s.sub(ctx, op); err != nil {
		return fmt.Errorf(`sub: %w`, err)
	}

	return nil
}

func (s *AccountStore) Lock(ctx router.Context, burn *BalanceOperation) (*LockId, error) {
	return nil, nil
}

func (s *AccountStore) Unlock(ctx router.Context, id *LockId) error {
	return nil
}

func (s *AccountStore) BurnLock(ctx router.Context, id *LockId) error {
	return nil
}

// todo: add TransferLock implementation
func (u *AccountStore) TransferLock(ctx router.Context, id *LockId, transfer *TransferOperation) error {
	return nil
}

func (s *AccountStore) BurnAllLock(router.Context, *BalanceOperation) error {
	return nil
}

func (s *AccountStore) LockAll(router.Context, *BalanceOperation) error {
	return nil
}

func (s *AccountStore) GetLocked(ctx router.Context, balanceId *BalanceId) (*Balance, error) {
	return nil, nil
}

func (s *AccountStore) add(ctx router.Context, op *BalanceOperation) (*Balance, error) {
	balance, err := s.Get(ctx, &BalanceId{Address: op.Address, Symbol: op.Symbol, Group: op.Group})
	if err != nil {
		return nil, err
	}

	curBalanceAmount, err := balance.Amount.BigInt()
	if err != nil {
		return nil, fmt.Errorf(`parse cur balance: %w`, err)
	}

	toAdd, err := op.Amount.BigInt()
	if err != nil {
		return nil, fmt.Errorf(`parse amnount to add: %w`, err)
	}
	newBalance := &Balance{
		Address: op.Address,
		Symbol:  op.Symbol,
		Group:   op.Group,
		Amount:  decimal.BigIntSubAsDecimal(curBalanceAmount, toAdd),
	}

	if err = State(ctx).Put(newBalance); err != nil {
		return newBalance, err
	}
	return newBalance, err
}

func (s *AccountStore) sub(ctx router.Context, op *BalanceOperation) (*Balance, error) {
	balance, err := s.Get(ctx, &BalanceId{Address: op.Address, Symbol: op.Symbol, Group: op.Group})
	if err != nil {
		return nil, err
	}

	balAmount, err := balance.Amount.BigInt()
	if err != nil {
		return nil, err
	}
	opAmount, err := op.Amount.BigInt()
	if err != nil {
		return nil, err
	}

	if balAmount.Cmp(opAmount) == -1 {
		return nil, fmt.Errorf(`subtract from=%s: %w`, op.Address, ErrAmountInsuficcient)
	}
	newBalance := &Balance{
		Address: op.Address,
		Symbol:  op.Symbol,
		Group:   op.Group,
		Amount:  decimal.BigIntSubAsDecimal(balAmount, opAmount),
	}

	if err = State(ctx).Put(newBalance); err != nil {
		return newBalance, err
	}
	return newBalance, err
}
