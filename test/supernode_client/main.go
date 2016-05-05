package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	gopath "path"
	"time"

	random "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-random"
	"github.com/ipfs/go-ipfs/thirdparty/ipfsaddr"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
	ma "gx/ipfs/QmcobAGsCjYt5DXoq9et9L8yR8er7o7Cu3DTvpaq12jYSz/go-multiaddr"

	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore"
	syncds "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-datastore/sync"
	commands "github.com/ipfs/go-ipfs/commands"
	core "github.com/ipfs/go-ipfs/core"
	corehttp "github.com/ipfs/go-ipfs/core/corehttp"
	corerouting "github.com/ipfs/go-ipfs/core/corerouting"
	"github.com/ipfs/go-ipfs/core/coreunix"
	"github.com/ipfs/go-ipfs/repo"
	config "github.com/ipfs/go-ipfs/repo/config"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	ds2 "github.com/ipfs/go-ipfs/thirdparty/datastore2"
	unit "github.com/ipfs/go-ipfs/thirdparty/unit"
	peer "gx/ipfs/QmZwZjMVGss5rqYsJVGy18gNbkTJffFyq2x1uJ4e4p3ZAt/go-libp2p-peer"
	logging "gx/ipfs/QmaDNZ4QMdBdku1YZWBysufYyoQt1negQGNav6PLYarbY8/go-log"
)

var elog = logging.Logger("gc-client")

var (
	cat             = flag.Bool("cat", false, "else add")
	seed            = flag.Int64("seed", 1, "")
	nBitsForKeypair = flag.Int("b", 1024, "number of bits for keypair (if repo is uninitialized)")
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	servers := config.DefaultSNRServers
	fmt.Println("using gcr remotes:")
	for _, p := range servers {
		fmt.Println("\t", p)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	repoPath := gopath.Join(cwd, config.DefaultPathName)
	if err := ensureRepoInitialized(repoPath); err != nil {
	}
	repo, err := fsrepo.Open(repoPath)
	if err != nil { // owned by node
		return err
	}
	cfg, err := repo.Config()
	if err != nil {
		return err
	}

	cfg.Bootstrap = servers
	if err := repo.SetConfig(cfg); err != nil {
		return err
	}

	var addrs []ipfsaddr.IPFSAddr
	for _, info := range servers {
		addr, err := ipfsaddr.ParseString(info)
		if err != nil {
			return err
		}
		addrs = append(addrs, addr)
	}

	var infos []peer.PeerInfo
	for _, addr := range addrs {
		infos = append(infos, peer.PeerInfo{
			ID:    addr.ID(),
			Addrs: []ma.Multiaddr{addr.Transport()},
		})
	}

	node, err := core.NewNode(ctx, &core.BuildCfg{
		Online:  true,
		Repo:    repo,
		Routing: corerouting.SupernodeClient(infos...),
	})
	if err != nil {
		return err
	}
	defer node.Close()

	opts := []corehttp.ServeOption{
		corehttp.CommandsOption(cmdCtx(node, repoPath)),
		corehttp.GatewayOption(false, nil),
	}

	if *cat {
		if err := runFileCattingWorker(ctx, node); err != nil {
			return err
		}
	} else {
		if err := runFileAddingWorker(node); err != nil {
			return err
		}
	}
	return corehttp.ListenAndServe(node, cfg.Addresses.API, opts...)
}

func ensureRepoInitialized(path string) error {
	if !fsrepo.IsInitialized(path) {
		conf, err := config.Init(ioutil.Discard, *nBitsForKeypair)
		if err != nil {
			return err
		}
		if err := fsrepo.Init(path, conf); err != nil {
			return err
		}
	}
	return nil
}

func sizeOfIthFile(i int64) int64 {
	return (1 << uint64(i)) * unit.KB
}

func runFileAddingWorker(n *core.IpfsNode) error {
	errs := make(chan error)
	go func() {
		var i int64
		for i = 1; i < math.MaxInt32; i++ {
			piper, pipew := io.Pipe()
			go func() {
				defer pipew.Close()
				if err := random.WritePseudoRandomBytes(sizeOfIthFile(i), pipew, *seed); err != nil {
					errs <- err
				}
			}()
			k, err := coreunix.Add(n, piper)
			if err != nil {
				errs <- err
			}
			log.Println("added file", "seed", *seed, "#", i, "key", k, "size", unit.Information(sizeOfIthFile(i)))
			time.Sleep(1 * time.Second)
		}
	}()

	var i int64
	for i = 0; i < math.MaxInt32; i++ {
		err := <-errs
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func runFileCattingWorker(ctx context.Context, n *core.IpfsNode) error {
	conf, err := config.Init(ioutil.Discard, *nBitsForKeypair)
	if err != nil {
		return err
	}

	r := &repo.Mock{
		D: ds2.CloserWrap(syncds.MutexWrap(datastore.NewMapDatastore())),
		C: *conf,
	}
	dummy, err := core.NewNode(ctx, &core.BuildCfg{
		Repo: r,
	})
	if err != nil {
		return err
	}

	errs := make(chan error)

	go func() {
		defer dummy.Close()
		var i int64 = 1
		for {
			buf := new(bytes.Buffer)
			if err := random.WritePseudoRandomBytes(sizeOfIthFile(i), buf, *seed); err != nil {
				errs <- err
			}
			// add to a dummy node to discover the key
			k, err := coreunix.Add(dummy, bytes.NewReader(buf.Bytes()))
			if err != nil {
				errs <- err
			}
			e := elog.EventBegin(ctx, "cat", logging.LoggableF(func() map[string]interface{} {
				return map[string]interface{}{
					"key":       k,
					"localPeer": n.Identity,
				}
			}))
			if r, err := coreunix.Cat(ctx, n, k); err != nil {
				e.Done()
				log.Printf("failed to cat file. seed: %d #%d key: %s err: %s", *seed, i, k, err)
			} else {
				log.Println("found file", "seed", *seed, "#", i, "key", k, "size", unit.Information(sizeOfIthFile(i)))
				io.Copy(ioutil.Discard, r)
				e.Done()
				log.Println("catted file", "seed", *seed, "#", i, "key", k, "size", unit.Information(sizeOfIthFile(i)))
				i++
			}
			time.Sleep(time.Second)
		}
	}()

	err = <-errs
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func toPeerInfos(bpeers []config.BootstrapPeer) ([]peer.PeerInfo, error) {
	var peers []peer.PeerInfo
	for _, bootstrap := range bpeers {
		p, err := toPeerInfo(bootstrap)
		if err != nil {
			return nil, err
		}
		peers = append(peers, p)
	}
	return peers, nil
}

func toPeerInfo(bootstrap config.BootstrapPeer) (p peer.PeerInfo, err error) {
	p = peer.PeerInfo{
		ID:    bootstrap.ID(),
		Addrs: []ma.Multiaddr{bootstrap.Multiaddr()},
	}
	return p, nil
}

func cmdCtx(node *core.IpfsNode, repoPath string) commands.Context {
	return commands.Context{
		Online:     true,
		ConfigRoot: repoPath,
		LoadConfig: func(path string) (*config.Config, error) {
			return node.Repo.Config()
		},
		ConstructNode: func() (*core.IpfsNode, error) {
			return node, nil
		},
	}
}
