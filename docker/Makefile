# Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
# Use of this source code is governed by Microsoft Reference Source
# License (MS-RSL) that can be found in the LICENSE file.

all:

link:
	@if [ -d "$(GOPATH)/src/github.com/ipfs/go-ipfs" ]; then echo "github.com/ipfs/go-ipfs package exists in GOPATH"; exit 1; fi
	ln -s $(shell pwd)/src $(GOPATH)/src/github.com/ipfs/go-ipfs

clean:
	rm -rf assets blocks blockservice cmd commands core exchange dagutils filestore \
	fuse go-* *.md LICENSE CODEOWNERS importer keystore merkledag mfs namesys p2p path pin \
	plugin repo tar thirdparty unixfs vendor quic-go

extract:
	# before must do: cd $$GOPATH/src/github.com/ipfs/go-ipfs && make install
	#
	surgical-extraction \
		--pkg github.com/ipfs/go-ipfs/cmd/ipfswatch \
		--out github.com/AtlantPlatform/go-ipfs \
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
	--unvendor go-reuseport
	# --rename badger:github.com/dgraph-io/badger
	
extract-apply:
	surgical-extraction \
		--pkg github.com/ipfs/go-ipfs/cmd/ipfswatch \
		--out github.com/AtlantPlatform/go-ipfs \
	extract --apply \
	--unvendor go-libp2p-pnet \
	--unvendor go-libp2p-interface-pnet \
	--unvendor go-libp2p-peer \
	--unvendor go-libp2p-crypto \
	--unvendor go-cid \
	--unvendor go-libp2p-kad-dht \
	--unvendor go-libp2p-secio \
	--unvendor go-libp2p-floodsub \
	--unvendor go-libp2p-peerstore \
	--unvendor go-block-format \
	--unvendor go-ipfs-cmdkit \
	--unvendor go-ipld-format \
	--unvendor go-multiaddr \
	--unvendor go-libp2p-metrics \
	--unvendor go-reuseport \
	--unvendor quic-go \
	--unvendor go-libp2p-transport \
	--unvendor go-libp2p-quic-transport \
	--unvendor go-libp2p
	rm -rf quic-go go-libp2p-quic-transport

# --unvendor go-ds-badger \
# --rename badger:github.com/dgraph-io/badger
replace:
	find . -type f | grep \.go$ | xargs sed -i \
	   's/github.com\/AtlantPlatform\/go-ipfs\/quic-go/github.com\/lucas-clemente\/quic-go/g'
	find . -type f | grep \.go$ | xargs sed -i \
	   's/github.com\/AtlantPlatform\/go-ipfs\/go-libp2p-quic-transport/github.com\/libp2p\/go-libp2p-quic-transport/g'

patch-apply:
	git apply patches/ed25519.patch
	git apply patches/libp2p-version.patch
	git apply patches/mdns.patch
	git apply patches/peer-info.patch
	# deprecated: git apply patches/go-reuseport.patch

test:
	go install github.com/AtlantPlatform/go-ipfs/cmd/ipfswatch
