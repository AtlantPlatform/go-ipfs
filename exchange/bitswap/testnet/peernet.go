package bitswap

import (
	"context"

	bsnet "bitbucket.org/atlantproject/go-ipfs/exchange/bitswap/network"
	ds "unknown/go-datastore"
	mockrouting "unknown/go-ipfs-routing/mock"
	peer "unknown/go-libp2p-peer"
	mockpeernet "unknown/go-libp2p/p2p/net/mock"
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
