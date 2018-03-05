package coredag

import (
	"io"
	"io/ioutil"

	ipldcbor "unknown/go-ipld-cbor"
	ipld "unknown/go-ipld-format"
)

func cborJSONParser(r io.Reader, mhType uint64, mhLen int) ([]ipld.Node, error) {
	nd, err := ipldcbor.FromJson(r, mhType, mhLen)
	if err != nil {
		return nil, err
	}

	return []ipld.Node{nd}, nil
}

func cborRawParser(r io.Reader, mhType uint64, mhLen int) ([]ipld.Node, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	nd, err := ipldcbor.Decode(data, mhType, mhLen)
	if err != nil {
		return nil, err
	}

	return []ipld.Node{nd}, nil
}
