package ifconnmgr

import (
	"context"
	"time"

	inet "unknown/go-libp2p-net"
	peer "unknown/go-libp2p-peer"
)

type ConnManager interface {
	TagPeer(peer.ID, string, int)
	UntagPeer(peer.ID, string)
	GetTagInfo(peer.ID) *TagInfo
	TrimOpenConns(context.Context)
	Notifee() inet.Notifiee
}

type TagInfo struct {
	FirstSeen time.Time
	Value     int
	Tags      map[string]int
	Conns     map[string]time.Time
}
