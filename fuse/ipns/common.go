package ipns

import (
	"context"

	"bitbucket.org/atlantproject/go-ipfs/core"
	ci "bitbucket.org/atlantproject/go-ipfs/go-libp2p-crypto"
	nsys "bitbucket.org/atlantproject/go-ipfs/namesys"
	path "bitbucket.org/atlantproject/go-ipfs/path"
	ft "bitbucket.org/atlantproject/go-ipfs/unixfs"
)

// InitializeKeyspace sets the ipns record for the given key to
// point to an empty directory.
func InitializeKeyspace(n *core.IpfsNode, key ci.PrivKey) error {
	ctx, cancel := context.WithCancel(n.Context())
	defer cancel()

	emptyDir := ft.EmptyDirNode()

	err := n.Pinning.Pin(ctx, emptyDir, false)
	if err != nil {
		return err
	}

	err = n.Pinning.Flush()
	if err != nil {
		return err
	}

	pub := nsys.NewRoutingPublisher(n.Routing, n.Repo.Datastore())

	return pub.Publish(ctx, key, path.FromCid(emptyDir.Cid()))
}
