package mockrouting

import (
	"context"
	"errors"
	"time"

	cid "github.com/AtlantPlatform/go-ipfs/go-cid"
	peer "github.com/AtlantPlatform/go-ipfs/go-libp2p-peer"
	pstore "github.com/AtlantPlatform/go-ipfs/go-libp2p-peerstore"
	ma "github.com/AtlantPlatform/go-ipfs/go-multiaddr"
	ds "unknown/go-datastore"
	dshelp "unknown/go-ipfs-ds-help"
	u "unknown/go-ipfs-util"
	dhtpb "unknown/go-libp2p-record/pb"
	routing "unknown/go-libp2p-routing"
	logging "unknown/go-log"
	"unknown/go-testutil"
	proto "unknown/gogo-protobuf/proto"
)

var log = logging.Logger("mockrouter")

type client struct {
	datastore ds.Datastore
	server    server
	peer      testutil.Identity
}

// FIXME(brian): is this method meant to simulate putting a value into the network?
func (c *client) PutValue(ctx context.Context, key string, val []byte) error {
	log.Debugf("PutValue: %s", key)
	rec := new(dhtpb.Record)
	rec.Value = val
	rec.Key = proto.String(string(key))
	rec.TimeReceived = proto.String(u.FormatRFC3339(time.Now()))
	data, err := proto.Marshal(rec)
	if err != nil {
		return err
	}

	return c.datastore.Put(dshelp.NewKeyFromBinary([]byte(key)), data)
}

// FIXME(brian): is this method meant to simulate getting a value from the network?
func (c *client) GetValue(ctx context.Context, key string) ([]byte, error) {
	log.Debugf("GetValue: %s", key)
	v, err := c.datastore.Get(dshelp.NewKeyFromBinary([]byte(key)))
	if err != nil {
		return nil, err
	}

	data, ok := v.([]byte)
	if !ok {
		return nil, errors.New("could not cast value from datastore")
	}

	rec := new(dhtpb.Record)
	err = proto.Unmarshal(data, rec)
	if err != nil {
		return nil, err
	}

	return rec.GetValue(), nil
}

func (c *client) GetValues(ctx context.Context, key string, count int) ([]routing.RecvdVal, error) {
	log.Debugf("GetValues: %s", key)
	data, err := c.GetValue(ctx, key)
	if err != nil {
		return nil, err
	}

	return []routing.RecvdVal{{Val: data, From: c.peer.ID()}}, nil
}

func (c *client) FindProviders(ctx context.Context, key *cid.Cid) ([]pstore.PeerInfo, error) {
	return c.server.Providers(key), nil
}

func (c *client) FindPeer(ctx context.Context, pid peer.ID) (pstore.PeerInfo, error) {
	log.Debugf("FindPeer: %s", pid)
	return pstore.PeerInfo{}, nil
}

func (c *client) FindProvidersAsync(ctx context.Context, k *cid.Cid, max int) <-chan pstore.PeerInfo {
	out := make(chan pstore.PeerInfo)
	go func() {
		defer close(out)
		for i, p := range c.server.Providers(k) {
			if max <= i {
				return
			}
			select {
			case out <- p:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Provide returns once the message is on the network. Value is not necessarily
// visible yet.
func (c *client) Provide(_ context.Context, key *cid.Cid, brd bool) error {
	if !brd {
		return nil
	}
	info := pstore.PeerInfo{
		ID:    c.peer.ID(),
		Addrs: []ma.Multiaddr{c.peer.Address()},
	}
	return c.server.Announce(info, key)
}

func (c *client) Ping(ctx context.Context, p peer.ID) (time.Duration, error) {
	return 0, nil
}

func (c *client) Bootstrap(context.Context) error {
	return nil
}

var _ routing.IpfsRouting = &client{}
