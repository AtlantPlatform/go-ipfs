package metrics

import (
	peer "github.com/AtlantPlatform/go-ipfs/go-libp2p-peer"
	protocol "unknown/go-libp2p-protocol"
)

type StreamMeterCallback func(int64, protocol.ID, peer.ID)
type MeterCallback func(int64)

type Reporter interface {
	LogSentMessage(int64)
	LogRecvMessage(int64)
	LogSentMessageStream(int64, protocol.ID, peer.ID)
	LogRecvMessageStream(int64, protocol.ID, peer.ID)
	GetBandwidthForPeer(peer.ID) Stats
	GetBandwidthForProtocol(protocol.ID) Stats
	GetBandwidthTotals() Stats
}
