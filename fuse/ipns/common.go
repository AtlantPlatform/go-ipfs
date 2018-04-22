package ipns

import (
	"context"

	"github.com/AtlantPlatform/go-ipfs/core"
	ci "github.com/AtlantPlatform/go-ipfs/go-libp2p-crypto"
	nsys "github.com/AtlantPlatform/go-ipfs/namesys"
	path "github.com/AtlantPlatform/go-ipfs/path"
	ft "github.com/AtlantPlatform/go-ipfs/unixfs"
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
