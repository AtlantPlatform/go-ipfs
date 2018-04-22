package bitswap

import (
	"context"

	bsnet "github.com/AtlantPlatform/go-ipfs/exchange/bitswap/network"
	peer "github.com/AtlantPlatform/go-ipfs/go-libp2p-peer"
	mockpeernet "github.com/AtlantPlatform/go-ipfs/go-libp2p/p2p/net/mock"
	ds "unknown/go-datastore"
	mockrouting "unknown/go-ipfs-routing/mock"
	testutil "unknown/go-testutil"
)

type peernet struct {
	mockpeernet.Mocknet
	routingserver mockrouting.Server
}

func StreamNet(ctx context.Context, net mockpeernet.Mocknet, rs mockrouting.Server) (Network, error) {
	return &peernet{net, rs}, nil
}

func (pn *peernet) Adapter(p testutil.Identity) bsnet.BitSwapNetwork {
	client, err := pn.Mocknet.AddPeer(p.PrivateKey(), p.Address())
	if err != nil {
		panic(err.Error())
	}
	routing := pn.routingserver.ClientWithDatastore(context.TODO(), p, ds.NewMapDatastore())
	return bsnet.NewFromIpfsHost(client, routing)
}

func (pn *peernet) HasPeer(p peer.ID) bool {
	for _, member := range pn.Mocknet.Peers() {
		if p == member {
			return true
		}
	}
	return false
}

var _ Network = (*peernet)(nil)
