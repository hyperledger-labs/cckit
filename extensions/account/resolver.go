package account

import (
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hyperledger-labs/cckit/router"
)

type Getter interface {
	GetInvokerAddress(router.Context, *emptypb.Empty) (*AddressId, error)

	GetAddress(router.Context, *GetAddressRequest) (*AddressId, error)

	GetAccount(router.Context, *AccountId) (*Account, error)
}
