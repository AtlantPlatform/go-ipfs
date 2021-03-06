// Package nilrouting implements a routing client that does nothing.
package nilrouting

import (
	"context"
	"errors"

	cid "github.com/AtlantPlatform/go-ipfs/go-cid"
	peer "github.com/AtlantPlatform/go-ipfs/go-libp2p-peer"
	pstore "github.com/AtlantPlatform/go-ipfs/go-libp2p-peerstore"
	ds "unknown/go-datastore"
	p2phost "unknown/go-libp2p-host"
	routing "unknown/go-libp2p-routing"
)

type nilclient struct {
}

func (c *nilclient) PutValue(_ context.Context, _ string, _ []byte) error {
	return nil
}

func (c *nilclient) GetValue(_ context.Context, _ string) ([]byte, error) {
	return nil, errors.New("tried GetValue from nil routing")
}

func (c *nilclient) GetValues(_ context.Context, _ string, _ int) ([]routing.RecvdVal, error) {
	return nil, errors.New("tried GetValues from nil routing")
}

func (c *nilclient) FindPeer(_ context.Context, _ peer.ID) (pstore.PeerInfo, error) {
	return pstore.PeerInfo{}, nil
}

func (c *nilclient) FindProvidersAsync(_ context.Context, _ *cid.Cid, _ int) <-chan pstore.PeerInfo {
	out := make(chan pstore.PeerInfo)
	defer close(out)
	return out
}

func (c *nilclient) Provide(_ context.Context, _ *cid.Cid, _ bool) error {
	return nil
}

func (c *nilclient) Bootstrap(_ context.Context) error {
	return nil
}

// ConstructNilRouting creates an IpfsRouting client which does nothing.
func ConstructNilRouting(_ context.Context, _ p2phost.Host, _ ds.Batching) (routing.IpfsRouting, error) {
	return &nilclient{}, nil
}

//  ensure nilclient satisfies interface
var _ routing.IpfsRouting = &nilclient{}
