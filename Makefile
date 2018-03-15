all:

link:
	@if [ -d "$(GOPATH)/src/github.com/ipfs/go-ipfs" ]; then echo "github.com/ipfs/go-ipfs package exists in GOPATH"; exit 1; fi
	ln -s $(shell pwd)/src $(GOPATH)/src/github.com/ipfs/go-ipfs

extract:
	# before must do: cd $$GOPATH/src/github.com/ipfs/go-ipfs && make install
	#
	surgical-extraction --pkg github.com/ipfs/go-ipfs/cmd/ipfswatch --out bitbucket.org/atlantproject/go-ipfs \
	extract \
	--unvendor go-libp2p-pnet \
	--unvendor go-libp2p-interface-pnet \
	--unvendor go-ds-badger \
	--unvendor go-libp2p-peer \
	--unvendor go-libp2p-crypto \
	--unvendor go-cid \
	--unvendor go-libp2p-kad-dht \
	--unvendor go-libp2p-secio \
	--unvendor go-libp2p \
	--rename badger:github.com/dgraph-io/badger
	
extract-apply:
	surgical-extraction --pkg github.com/ipfs/go-ipfs/cmd/ipfswatch --out bitbucket.org/atlantproject/go-ipfs \
	extract --apply \
	--unvendor go-libp2p-pnet \
	--unvendor go-libp2p-interface-pnet \
	--unvendor go-ds-badger \
	--unvendor go-libp2p-peer \
	--unvendor go-libp2p-crypto \
	--unvendor go-cid \
	--unvendor go-libp2p-kad-dht \
	--unvendor go-libp2p-secio \
	--unvendor go-libp2p \
	--rename badger:github.com/dgraph-io/badger

patch-apply:
	git apply patches/ed25519.patch

test:
	go install bitbucket.org/atlantproject/go-ipfs/cmd/ipfswatch
