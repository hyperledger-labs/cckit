package token

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/hyperledger-labs/cckit/extensions/token/decimal"
	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/state"
)

var _ BalanceStore = &UTXOStore{}

var (
	ErrSenderRecipientEqual = errors.New(`sender recipient equal`)
	ErrSenderNotEqual       = errors.New(`sender not equal`)
	ErrOpAddressNotEqual    = errors.New(`operation address not equal`)
	ErrSymbolNotEqual       = errors.New(`symbol not equal`)
	ErrRecipientDuplicate   = errors.New(`errors recipient duplicate`)
)

type UTXOStore struct {
}

func (u *UTXO) ID() *UTXOId {
	return &UTXOId{
		Address: u.Address,
		Symbol:  u.Symbol,
		Group:   u.Group,
		TxId:    u.TxId,
	}
}

func NewUTXOStore() *UTXOStore {
	return &UTXOStore{}
}

func (u *UTXOStore) Get(ctx router.Context, balanceId *BalanceId) (*Balance, error) {
	if err := router.ValidateRequest(balanceId); err != nil {
		return nil, err
	}
	outputs, err := listOutputs(ctx, balanceId)
	if err != nil {
		return nil, err
	}

	amount := big.NewInt(0)
	for _, u := range outputs {
		uAmount, err := u.Amount.BigInt()
		if err != nil {
			return nil, err
		}
		amount.Add(amount, uAmount)
	}

	if amount == nil {
		amount = big.NewInt(0)
	}
	balance := &Balance{
		Address: balanceId.Address,
		Symbol:  balanceId.Symbol,
		Group:   balanceId.Group,
		Amount:  decimal.New(amount),
	}

	return balance, nil
}

func (u *UTXOStore) GetLocked(ctx router.Context, balanceId *BalanceId) (*Balance, error) {
	if err := router.ValidateRequest(balanceId); err != nil {
		return nil, err
	}
	outputs, err := listOutputs(ctx, balanceId)
	if err != nil {
		return nil, err
	}

	amount := big.NewInt(0)
	for _, u := range outputs {
		if u.Locked {
			uAmount, err := u.Amount.BigInt()
			if err != nil {
				return nil, err
			}
			amount.Add(amount, uAmount)
		}
	}

	balance := &Balance{
		Symbol:  balanceId.Symbol,
		Group:   balanceId.Group,
		Address: balanceId.Address,
		Amount:  decimal.New(amount),
	}

	return balance, nil
}

func (u *UTXOStore) GetUTXO(ctx router.Context, utxoId *UTXOId) (*UTXO, error) {
	utxo, err := State(ctx).Get(&UTXOId{Symbol: utxoId.Symbol, Group: utxoId.Group, Address: utxoId.Address, TxId: utxoId.TxId}, &UTXO{})
	if err != nil {
		return nil, err
	}
	return utxo.(*UTXO), nil
}

// ListOutputs unspended outputs list
func (u *UTXOStore) ListOutputs(ctx router.Context, balanceId *BalanceId) ([]*UTXO, error) {
	return listOutputs(ctx, balanceId)
}

func (u *UTXOStore) List(ctx router.Context, id *TokenId) ([]*Balance, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UTXOStore) Transfer(ctx router.Context, transfer *TransferOperation) error {
	if err := router.ValidateRequest(transfer); err != nil {
		return err
	}

	if transfer.Sender == transfer.Recipient {
		return ErrSenderRecipientEqual
	}

	senderOutputs, err := u.ListOutputs(ctx, &BalanceId{
		Symbol:  transfer.Symbol,
		Group:   transfer.Group,
		Address: transfer.Sender,
	})
	if err != nil {
		return err
	}

	transferAmount, err := transfer.Amount.BigInt()
	if err != nil {
		return err
	}
	useOutputs, outputsAmount, err := selectOutputsForAmount(senderOutputs, transferAmount, false)
	if err != nil {
		return err
	}

	txID := ctx.Stub().GetTxID() + ".0"
	recipientOutput := &UTXO{
		Symbol:  transfer.Symbol,
		Group:   strings.Join(transfer.Group, `,`),
		Address: transfer.Recipient,
		TxId:    txID,
		Amount:  transfer.Amount,
		Locked:  false,
		//Meta: transfer.Meta,
	}

	if err := State(ctx).Insert(recipientOutput); err != nil {
		return err
	}

	if outputsAmount.Cmp(transferAmount) == 1 {
		senderChangeOutput := &UTXO{
			Symbol:  transfer.Symbol,
			Group:   strings.Join(transfer.Group, `,`),
			Address: transfer.Sender,
			TxId:    txID,
			Amount:  decimal.BigIntSubAsDecimal(outputsAmount, transferAmount, transfer.Amount.Scale),
			Locked:  false,
		}
		if err := State(ctx).Insert(senderChangeOutput); err != nil {
			return err
		}
	}

	for _, output := range useOutputs {
		if err := State(ctx).Delete(output.ID()); err != nil {
			return err
		}
	}

	return nil
}

func (u *UTXOStore) TransferBatch(ctx router.Context, transfers []*TransferOperation) error {
	var (
		sender, symbol string
		scale          int32
		group          []string
		recipients     = make(map[string]interface{})
		totalAmount    = big.NewInt(0)
	)
	for _, transfer := range transfers {

		if err := router.ValidateRequest(transfer); err != nil {
			return err
		}

		if sender == `` {
			sender = transfer.Sender
		}

		if transfer.Sender != sender {
			return ErrSenderNotEqual
		}

		if sender == transfer.Recipient {
			return ErrSenderRecipientEqual
		}
		if symbol == `` {
			symbol = transfer.Symbol
		}

		if transfer.Symbol != symbol {
			return ErrSymbolNotEqual
		}

		if len(transfer.Group) > 0 {
			panic(`implement me`)
		}

		scale = transfer.Amount.Scale

		if _, ok := recipients[transfer.Recipient]; ok {
			return ErrRecipientDuplicate
		}
		recipients[transfer.Recipient] = nil
		transferAmount, err := transfer.Amount.BigInt()
		if err != nil {
			return err
		}
		totalAmount.Add(totalAmount, transferAmount)
		// totalAmount += transfer.Amount
	}

	senderOutputs, err := u.ListOutputs(ctx, &BalanceId{
		Symbol:  symbol,
		Group:   group,
		Address: sender,
	})
	if err != nil {
		return err
	}

	useOutputs, outputsAmount, err := selectOutputsForAmount(senderOutputs, totalAmount, false)
	if err != nil {
		return err
	}

	for _, output := range useOutputs {
		if err := State(ctx).Delete(output.ID()); err != nil {
			return err
		}
	}

	txID := ctx.Stub().GetTxID() + ".0"

	for _, transfer := range transfers {
		recipientOutput := &UTXO{
			Symbol:  transfer.Symbol,          // INV
			Group:   strings.Join(group, `,`), // 001
			Address: transfer.Recipient,
			TxId:    txID,
			Amount:  transfer.Amount,
			Locked:  false,
			//Meta: transfer.Meta,
		}
		if err := State(ctx).Insert(recipientOutput); err != nil {
			return err
		}
	}

	if outputsAmount.Cmp(totalAmount) == 1 {
		senderChangeOutput := &UTXO{
			Symbol:  symbol,
			Group:   strings.Join(group, `,`),
			Address: sender,
			TxId:    txID,
			Amount:  decimal.BigIntSubAsDecimal(outputsAmount, totalAmount, scale),
			Locked:  false,
		}
		if err := State(ctx).Insert(senderChangeOutput); err != nil {
			return err
		}
	}

	return nil
}

func (u *UTXOStore) Mint(ctx router.Context, op *BalanceOperation) error {
	mintedOutput := &UTXO{
		Address: op.Address,
		Symbol:  op.Symbol,
		Group:   strings.Join(op.Group, `,`),
		TxId:    ctx.Stub().GetTxID() + ".0",
		Amount:  op.Amount,
		Locked:  false,
	}

	return State(ctx).Insert(mintedOutput)
}

// Lock tokens
func (u *UTXOStore) Lock(ctx router.Context, op *BalanceOperation) (*LockId, error) {
	if err := router.ValidateRequest(op); err != nil {
		return nil, err
	}

	outputs, err := u.ListOutputs(ctx, &BalanceId{
		Symbol:  op.Symbol,
		Group:   op.Group,
		Address: op.Address,
	})
	if err != nil {
		return nil, err
	}

	opAmount, err := op.Amount.BigInt()
	if err != nil {
		return nil, err
	}
	useOutputs, outputsAmount, err := selectOutputsForAmount(outputs, opAmount, false)
	if err != nil {
		return nil, err
	}

	for _, output := range useOutputs {
		if err := State(ctx).Delete(output.ID()); err != nil {
			return nil, err
		}
	}

	lockedOutput := &UTXO{
		Symbol:  op.Symbol,
		Group:   strings.Join(op.Group, `,`),
		Address: op.Address,
		TxId:    ctx.Stub().GetTxID() + ".0",
		Amount:  op.Amount,
		Locked:  true,
		//Meta: transfer.Meta,
	}

	if err := State(ctx).Insert(lockedOutput); err != nil {
		return nil, err
	}

	if outputsAmount.Cmp(opAmount) == 1 {
		senderChangeOutput := &UTXO{
			Symbol:  op.Symbol,
			Group:   strings.Join(op.Group, `,`),
			Address: op.Address,
			TxId:    ctx.Stub().GetTxID() + ".1",
			Amount:  decimal.BigIntSubAsDecimal(outputsAmount, opAmount, op.Amount.Scale),
			Locked:  false,
		}
		if err := State(ctx).Insert(senderChangeOutput); err != nil {
			return nil, err
		}
	}

	return &LockId{lockedOutput.Symbol, lockedOutput.Group, lockedOutput.Address, lockedOutput.TxId}, nil
}

func (u *UTXOStore) LockBatch(ctx router.Context, ops []*BalanceOperation) ([]*LockId, error) {
	var (
		opAddress, symbol string
		scale             int32
		group             []string
		totalAmount       = big.NewInt(0)
	)

	for _, op := range ops {
		if err := router.ValidateRequest(op); err != nil {
			return nil, err
		}

		if opAddress == `` {
			opAddress = op.Address
		}
		if op.Address != opAddress {
			return nil, ErrOpAddressNotEqual
		}

		if symbol == `` {
			symbol = op.Symbol
		}
		if op.Symbol != symbol {
			return nil, ErrSymbolNotEqual
		}

		if len(op.Group) > 0 {
			panic(`implement me`)
		}

		scale = op.Amount.Scale

		opAmount, err := op.Amount.BigInt()
		if err != nil {
			return nil, err
		}
		totalAmount.Add(totalAmount, opAmount)
	}

	outputs, err := u.ListOutputs(ctx, &BalanceId{
		Symbol:  symbol,
		Group:   group,
		Address: opAddress,
	})
	if err != nil {
		return nil, err
	}

	useOutputs, outputsAmount, err := selectOutputsForAmount(outputs, totalAmount, false)
	if err != nil {
		return nil, err
	}

	for _, output := range useOutputs {
		if err := State(ctx).Delete(output.ID()); err != nil {
			return nil, err
		}
	}

	lockIds := make([]*LockId, len(ops))

	for id, op := range ops {
		txID := ctx.Stub().GetTxID() + fmt.Sprintf(".%v", id)

		lockedOutput := &UTXO{
			Symbol:  op.Symbol,
			Group:   strings.Join(op.Group, `,`),
			Address: op.Address,
			TxId:    txID,
			Amount:  op.Amount,
			Locked:  true,
			//Meta: transfer.Meta,
		}

		lockIds[id] = &LockId{
			Symbol:  op.Symbol,
			Group:   strings.Join(op.Group, `,`),
			Address: op.Address,
			TxId:    txID,
		}

		if err := State(ctx).Insert(lockedOutput); err != nil {
			return nil, err
		}
	}

	if outputsAmount.Cmp(totalAmount) == 1 {
		senderChangeOutput := &UTXO{
			Symbol:  symbol,
			Group:   strings.Join(group, `,`),
			Address: opAddress,
			TxId:    ctx.Stub().GetTxID() + fmt.Sprintf(".%v", len(ops)),
			Amount:  decimal.BigIntSubAsDecimal(outputsAmount, totalAmount, scale),
			Locked:  false,
		}
		if err := State(ctx).Insert(senderChangeOutput); err != nil {
			return nil, err
		}
	}

	return lockIds, nil
}

func (u *UTXOStore) LockAll(ctx router.Context, op *BalanceOperation) error {
	utxos, err := State(ctx).ListWith(&UTXO{}, state.Key{op.Symbol, strings.Join(op.Group, `,`)}) // todo: ???
	if err != nil {
		return err
	}

	for _, output := range utxos.(*UTXOs).Items {
		output.Locked = true
		if err := State(ctx).Put(output); err != nil {
			return err
		}
	}
	return nil
}

// Unlock tokens
func (u *UTXOStore) Unlock(ctx router.Context, id *LockId) error {
	utxo, err := State(ctx).Get(&UTXOId{Symbol: id.Symbol, Group: id.Group, Address: id.Address, TxId: id.TxId}, &UTXO{})
	if err != nil {
		return err
	}
	lockedOutput := utxo.(*UTXO)
	lockedOutput.Locked = false

	if err := State(ctx).Put(lockedOutput); err != nil {
		return err
	}
	return nil
}

// Burn unlocked tokens
func (u *UTXOStore) Burn(ctx router.Context, op *BalanceOperation) error {
	return burn(ctx, op, false)
}

// BurnLock locked tokens
func (u *UTXOStore) BurnLock(ctx router.Context, id *LockId) error {
	utxo, err := State(ctx).Get(&UTXOId{Symbol: id.Symbol, Group: id.Group, Address: id.Address, TxId: id.TxId}, &UTXO{})
	if err != nil {
		return err
	}
	lockedOutput := utxo.(*UTXO)

	if err := State(ctx).Delete(lockedOutput.ID()); err != nil {
		return err
	}
	return nil
}

// TransferLock locked tokens between accounts
func (u *UTXOStore) TransferLock(ctx router.Context, id *LockId, transfer *TransferOperation) error {
	utxo, err := State(ctx).Get(&UTXOId{Symbol: id.Symbol, Group: id.Group, Address: id.Address, TxId: id.TxId}, &UTXO{})
	if err != nil {
		return err
	}
	transferLockedOutput := utxo.(*UTXO)
	if err := State(ctx).Delete(transferLockedOutput.ID()); err != nil {
		return err
	}

	transferLockedOutput.Address = transfer.Recipient
	transferLockedOutput.Locked = false // unlock token
	if err := State(ctx).Put(transferLockedOutput); err != nil {
		return err
	}

	return nil
}

func (u *UTXOStore) BurnAllLock(ctx router.Context, op *BalanceOperation) error {
	utxos, err := State(ctx).ListWith(&UTXO{}, state.Key{op.Symbol, strings.Join(op.Group, `,`)}) // todo: ???
	if err != nil {
		return err
	}

	for _, output := range utxos.(*UTXOs).Items {
		if output.Locked {
			if err := State(ctx).Delete(output.ID()); err != nil {
				return err
			}
		}
	}
	return nil
}

// todo: Optimize selection, to maximum fit outputs
func selectOutputsForAmount(outputs []*UTXO, amount *big.Int, locked bool) ([]*UTXO, *big.Int, error) {
	var (
		selectedOutputs []*UTXO
		curAmount       = big.NewInt(0)
	)

	for _, o := range outputs {
		if (!locked && !o.Locked) || (locked && o.Locked) {
			selectedOutputs = append(selectedOutputs, o)
			oAmount, err := o.Amount.BigInt()
			if err != nil {
				return nil, new(big.Int), err
			}
			curAmount.Add(curAmount, oAmount)
			if curAmount.Cmp(amount) >= 0 {
				return selectedOutputs, curAmount, nil
			}
		}
	}

	return nil, new(big.Int), ErrAmountInsuficcient
}

// get unspended outputs list
func listOutputs(ctx router.Context, balanceId *BalanceId) ([]*UTXO, error) {
	utxos, err := State(ctx).ListWith(&UTXO{}, UTXOKeyBase(&UTXO{
		Symbol:  balanceId.Symbol,
		Group:   strings.Join(balanceId.Group, `,`),
		Address: balanceId.Address,
	}))
	if err != nil {
		return nil, err
	}

	return utxos.(*UTXOs).Items, nil
}

// burn locked or unlocked tokens
func burn(ctx router.Context, burn *BalanceOperation, locked bool) error {
	outputs, err := listOutputs(ctx, &BalanceId{
		Address: burn.Address,
		Symbol:  burn.Symbol,
		Group:   burn.Group,
	})
	if err != nil {
		return err
	}

	burnAmount, err := burn.Amount.BigInt()
	if err != nil {
		return err
	}
	useOutputs, outputsAmount, err := selectOutputsForAmount(outputs, burnAmount, locked)
	if err != nil {
		return err
	}

	for _, output := range useOutputs {
		if err := State(ctx).Delete(output.ID()); err != nil {
			return err
		}
	}

	if outputsAmount.Cmp(burnAmount) == 1 {

		senderChangeOutput := &UTXO{
			Symbol:  burn.Symbol,
			Group:   strings.Join(burn.Group, `,`),
			Address: burn.Address,
			TxId:    ctx.Stub().GetTxID(),
			Amount:  decimal.BigIntSubAsDecimal(outputsAmount, burnAmount, burn.Amount.Scale),
			Locked:  locked,
		}
		if err := State(ctx).Insert(senderChangeOutput); err != nil {
			return err
		}
	}

	return nil
}

// Lock or unlock tokens
/*
func setLock(ctx router.Context, op *BalanceOperation, locked bool) error {
	outputs, err := listOutputs(ctx, &BalanceId{
		Address: op.Address,
		Symbol:  op.Symbol,
		Group:   op.Group,
	})
	if err != nil {
		return err
	}

	opAmount, err := utils.IntVal(op.Amount)
	if err != nil {
		return err
	}
	useOutputs, outputsAmount, err := selectOutputsForAmount(outputs, opAmount, !locked)
	if err != nil {
		return err
	}

	for _, output := range useOutputs {
		if err := State(ctx).Delete(output.ID()); err != nil {
			return err
		}
	}

	lockedOutput := &UTXO{
		Address: op.Address,
		Symbol:  op.Symbol,
		Group:   strings.Join(op.Group, `,`),
		TxId:    ctx.Stub().GetTxID(),
		Amount:  op.Amount,
		Locked:  locked,
	}
	if err := State(ctx).Insert(lockedOutput); err != nil {
		return err
	}

	if outputsAmount.Cmp(opAmount) == 1 {
		changeAmount := outputsAmount.Sub(outputsAmount, opAmount)
		senderChangeOutput := &UTXO{
			Address: op.Address,
			Symbol:  op.Symbol,
			Group:   strings.Join(op.Group, `,`),
			TxId:    ctx.Stub().GetTxID(),
			Amount:  changeNewBigInt(amount)
			Locked:  !locked,
		}
		if err := State(ctx).Insert(senderChangeOutput); err != nil {
			return err
		}
	}

	return nil
}
*/
