// Package importer implements utilities used to create IPFS DAGs from files
// and readers.
package importer

import (
	"fmt"
	"os"

	chunker "unknown/go-ipfs-chunker"
	"unknown/go-ipfs-cmdkit/files"
	ipld "unknown/go-ipld-format"

	bal "bitbucket.org/atlantproject/go-ipfs/importer/balanced"
	h "bitbucket.org/atlantproject/go-ipfs/importer/helpers"
	trickle "bitbucket.org/atlantproject/go-ipfs/importer/trickle"
)

// BuildDagFromFile builds a DAG from the given file, writing created blocks to
// disk as they are created.
func BuildDagFromFile(fpath string, ds ipld.DAGService) (ipld.Node, error) {
	stat, err := os.Lstat(fpath)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("`%s` is a directory", fpath)
	}

	f, err := files.NewSerialFile(fpath, fpath, false, stat)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return BuildDagFromReader(ds, chunker.DefaultSplitter(f))
}

// BuildDagFromReader creates a DAG given a DAGService and a Splitter
// implementation (Splitters are io.Readers), using a Balanced layout.
func BuildDagFromReader(ds ipld.DAGService, spl chunker.Splitter) (ipld.Node, error) {
	dbp := h.DagBuilderParams{
		Dagserv:  ds,
		Maxlinks: h.DefaultLinksPerBlock,
	}

	return bal.Layout(dbp.New(spl))
}

// BuildTrickleDagFromReader creates a DAG given a DAGService and a Splitter
// implementation (Splitters are io.Readers), using a Trickle Layout.
func BuildTrickleDagFromReader(ds ipld.DAGService, spl chunker.Splitter) (ipld.Node, error) {
	dbp := h.DagBuilderParams{
		Dagserv:  ds,
		Maxlinks: h.DefaultLinksPerBlock,
	}

	return trickle.Layout(dbp.New(spl))
}
