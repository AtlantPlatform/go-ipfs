# Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
# Use of this source code is governed by Microsoft Reference Source
# License (MS-RSL) that can be found in the LICENSE file.

all:

link:
	@if [ -d "$(GOPATH)/src/github.com/ipfs/go-ipfs" ]; then echo "github.com/ipfs/go-ipfs package exists in GOPATH"; exit 1; fi
	ln -s $(shell pwd)/src $(GOPATH)/src/github.com/ipfs/go-ipfs

extract:
	# before must do: cd $$GOPATH/src/github.com/ipfs/go-ipfs && make install
	#
	surgical-extraction --pkg github.com/ipfs/go-ipfs/cmd/ipfswatch --out github.com/AtlantPlatform/go-ipfs \
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
	--unvendor go-libp2p-floodsub \
	--unvendor go-libp2p-peerstore \
	--unvendor go-block-format \
	--unvendor go-ipfs-cmdkit \
	--unvendor go-ipld-format \
	--unvendor go-multiaddr \
	--unvendor go-libp2p-metrics \
	--unvendor go-reuseport \
	--rename badger:github.com/dgraph-io/badger
	
extract-apply:
	surgical-extraction --pkg github.com/ipfs/go-ipfs/cmd/ipfswatch --out github.com/AtlantPlatform/go-ipfs \
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
	--unvendor go-libp2p-floodsub \
	--unvendor go-libp2p-peerstore \
	--unvendor go-block-format \
	--unvendor go-ipfs-cmdkit \
	--unvendor go-ipld-format \
	--unvendor go-multiaddr \
	--unvendor go-libp2p-metrics \
	--unvendor go-reuseport \
	--rename badger:github.com/dgraph-io/badger

patch-apply:
	git apply patches/ed25519.patch
	git apply patches/libp2p-version.patch
	git apply patches/mdns.patch
	git apply patches/peer-info.patch
	git apply patches/go-reuseport.patch

test:
	go install github.com/AtlantPlatform/go-ipfs/cmd/ipfswatch
