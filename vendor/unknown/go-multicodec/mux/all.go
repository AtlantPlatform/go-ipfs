package muxcodec

import (
	mc "unknown/go-multicodec"
	cbor "unknown/go-multicodec/cbor"
	json "unknown/go-multicodec/json"
)

func StandardMux() *Multicodec {
	return MuxMulticodec([]mc.Multicodec{
		cbor.Multicodec(),
		json.Multicodec(false),
		json.Multicodec(true),
	}, SelectFirst)
}
