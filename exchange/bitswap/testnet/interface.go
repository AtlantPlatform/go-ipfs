package bitswap

import (
	bsnet "github.com/AtlantPlatform/go-ipfs/exchange/bitswap/network"
	peer "github.com/AtlantPlatform/go-ipfs/go-libp2p-peer"
	"unknown/go-testutil"
)

type Network interface {
	Adapter(testutil.Identity) bsnet.BitSwapNetwork

	HasPeer(peer.ID) bool
}
