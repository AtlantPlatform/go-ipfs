package coreunix

import (
	"context"

	core "bitbucket.org/atlantproject/go-ipfs/core"
	path "bitbucket.org/atlantproject/go-ipfs/path"
	resolver "bitbucket.org/atlantproject/go-ipfs/path/resolver"
	uio "bitbucket.org/atlantproject/go-ipfs/unixfs/io"
)

func Cat(ctx context.Context, n *core.IpfsNode, pstr string) (uio.DagReader, error) {
	r := &resolver.Resolver{
		DAG:         n.DAG,
		ResolveOnce: uio.ResolveUnixfsOnce,
	}

	dagNode, err := core.Resolve(ctx, n.Namesys, r, path.Path(pstr))
	if err != nil {
		return nil, err
	}

	return uio.NewDagReader(ctx, dagNode, n.DAG)
}
