package router

import "github.com/hyperledger-labs/cckit/serialize"

type RouterOpt func(*Group)

func WithSerializer(s serialize.Serializer) RouterOpt {
	return func(g *Group) {
		g.serializer = s
	}
}
