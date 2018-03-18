package mdutils

import (
	bsrv "bitbucket.org/atlantproject/go-ipfs/blockservice"
	"bitbucket.org/atlantproject/go-ipfs/exchange/offline"
	ipld "bitbucket.org/atlantproject/go-ipfs/go-ipld-format"
	dag "bitbucket.org/atlantproject/go-ipfs/merkledag"
	ds "unknown/go-datastore"
	dssync "unknown/go-datastore/sync"
	blockstore "unknown/go-ipfs-blockstore"
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
