package pinger

import (
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/hyperledger-labs/cckit/router"
)

func NewService() *ChaincodePinger {
	return &ChaincodePinger{}
}

type ChaincodePinger struct{}

func (c *ChaincodePinger) Ping(ctx router.Context, _ *empty.Empty) (*PingInfo, error) {
	i, err := Ping(ctx)
	if err != nil {
		return nil, err
	}

	return i.(*PingInfo), err
}
