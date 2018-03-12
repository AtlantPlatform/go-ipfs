package bitswap

import (
	bsnet "bitbucket.org/atlantproject/go-ipfs/exchange/bitswap/network"
	peer "bitbucket.org/atlantproject/go-ipfs/go-libp2p-peer"
	"unknown/go-testutil"
)

type Network interface {
	Adapter(testutil.Identity) bsnet.BitSwapNetwork

	HasPeer(peer.ID) bool
}
