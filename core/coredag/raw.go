package coredag

import (
	"io"
	"io/ioutil"
	"math"

	block "github.com/AtlantPlatform/go-ipfs/go-block-format"
	cid "github.com/AtlantPlatform/go-ipfs/go-cid"
	ipld "github.com/AtlantPlatform/go-ipfs/go-ipld-format"
	"github.com/AtlantPlatform/go-ipfs/merkledag"
	mh "unknown/go-multihash"
)

func rawRawParser(r io.Reader, mhType uint64, mhLen int) ([]ipld.Node, error) {
	if mhType == math.MaxUint64 {
		mhType = mh.SHA2_256
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	h, err := mh.Sum(data, mhType, mhLen)
	if err != nil {
		return nil, err
	}
	c := cid.NewCidV1(cid.Raw, h)
	blk, err := block.NewBlockWithCid(data, c)
	if err != nil {
		return nil, err
	}
	nd := &merkledag.RawNode{Block: blk}
	return []ipld.Node{nd}, nil
}
