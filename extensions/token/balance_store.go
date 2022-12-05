package token

import (
	"github.com/hyperledger-labs/cckit/router"
)

type (
	LockId struct {
		Symbol  string
		Group   string
		Address string
		TxId    string
	}

	Store interface {
		Get(router.Context, *BalanceId) (*Balance, error)
		GetLocked(router.Context, *BalanceId) (*Balance, error)
		List(router.Context, *TokenId) ([]*Balance, error)
		Transfer(router.Context, *TransferOperation) error
		TransferBatch(router.Context, []*TransferOperation) error
		Mint(router.Context, *BalanceOperation) error
		Burn(router.Context, *BalanceOperation) error
		Lock(router.Context, *BalanceOperation) (*LockId, error)
		Unlock(router.Context, *LockId) error
		BurnLock(router.Context, *LockId) error
		LockAll(router.Context, *BalanceOperation) error
		BurnAllLock(router.Context, *BalanceOperation) error
	}
)
