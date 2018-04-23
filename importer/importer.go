// Package importer implements utilities used to create IPFS DAGs from files
// and readers.
package importer

import (
	ipld "github.com/AtlantPlatform/go-ipfs/go-ipld-format"
	chunker "unknown/go-ipfs-chunker"

	bal "github.com/AtlantPlatform/go-ipfs/importer/balanced"
	h "github.com/AtlantPlatform/go-ipfs/importer/helpers"
	trickle "github.com/AtlantPlatform/go-ipfs/importer/trickle"
)

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
