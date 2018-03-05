package basemux

import (
	mc "unknown/go-multicodec"
	mux "unknown/go-multicodec/mux"

	b64 "unknown/go-multicodec/base/b64"
	bin "unknown/go-multicodec/base/bin"
	hex "unknown/go-multicodec/base/hex"
)

func AllBasesMux() *mux.Multicodec {
	m := mux.MuxMulticodec([]mc.Multicodec{
		hex.Multicodec(),
		b64.Multicodec(),
		bin.Multicodec(),
	}, mux.SelectFirst)
	m.Wrap = false
	return m
}
