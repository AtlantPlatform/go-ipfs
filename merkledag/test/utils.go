package mdutils

import (
	bsrv "github.com/AtlantPlatform/go-ipfs/blockservice"
	ipld "github.com/AtlantPlatform/go-ipfs/go-ipld-format"
	dag "github.com/AtlantPlatform/go-ipfs/merkledag"
	ds "unknown/go-datastore"
	dssync "unknown/go-datastore/sync"
	blockstore "unknown/go-ipfs-blockstore"
	offline "unknown/go-ipfs-exchange-offline"
)

// Mock returns a new thread-safe, mock DAGService.
func Mock() ipld.DAGService {
	return dag.NewDAGService(Bserv())
}

// Bserv returns a new, thread-safe, mock BlockService.
func Bserv() bsrv.BlockService {
	bstore := blockstore.NewBlockstore(dssync.MutexWrap(ds.NewMapDatastore()))
	return bsrv.New(bstore, offline.Exchange(bstore))
}
