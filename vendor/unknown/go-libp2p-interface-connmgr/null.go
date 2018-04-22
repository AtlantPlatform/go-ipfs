package ifconnmgr

import (
	"context"

	peer "github.com/AtlantPlatform/go-ipfs/go-libp2p-peer"
	ma "github.com/AtlantPlatform/go-ipfs/go-multiaddr"
	inet "unknown/go-libp2p-net"
)

type NullConnMgr struct{}

func (_ NullConnMgr) TagPeer(peer.ID, string, int)  {}
func (_ NullConnMgr) UntagPeer(peer.ID, string)     {}
func (_ NullConnMgr) GetTagInfo(peer.ID) *TagInfo   { return &TagInfo{} }
func (_ NullConnMgr) TrimOpenConns(context.Context) {}
func (_ NullConnMgr) Notifee() inet.Notifiee        { return &cmNotifee{} }

type cmNotifee struct{}

func (nn *cmNotifee) Connected(n inet.Network, c inet.Conn)         {}
func (nn *cmNotifee) Disconnected(n inet.Network, c inet.Conn)      {}
func (nn *cmNotifee) Listen(n inet.Network, addr ma.Multiaddr)      {}
func (nn *cmNotifee) ListenClose(n inet.Network, addr ma.Multiaddr) {}
func (nn *cmNotifee) OpenedStream(inet.Network, inet.Stream)        {}
func (nn *cmNotifee) ClosedStream(inet.Network, inet.Stream)        {}
