package websocket

import (
	"context"

	ma "bitbucket.org/atlantproject/go-ipfs/go-multiaddr"
	tpt "unknown/go-libp2p-transport"
	manet "unknown/go-multiaddr-net"
	ws "unknown/websocket"
)

type dialer struct{}

func (d *dialer) Dial(raddr ma.Multiaddr) (tpt.Conn, error) {
	return d.DialContext(context.Background(), raddr)
}

func (d *dialer) DialContext(ctx context.Context, raddr ma.Multiaddr) (tpt.Conn, error) {
	wsurl, err := parseMultiaddr(raddr)
	if err != nil {
		return nil, err
	}

	wscon, _, err := ws.DefaultDialer.Dial(wsurl, nil)
	if err != nil {
		return nil, err
	}

	mnc, err := manet.WrapNetConn(NewConn(wscon, nil))
	if err != nil {
		wscon.Close()
		return nil, err
	}

	return &wsConn{
		Conn: mnc,
	}, nil
}

func (d *dialer) Matches(a ma.Multiaddr) bool {
	return WsFmt.Matches(a)
}
