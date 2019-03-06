FROM golang:1.11.5-stretch
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get -y update \
    && apt-get install -y git gcc libc-dev make

# - - - - - -
# "Badger" dependency has no backward compatibility,
# So we are taking it from the old branch
# and it has own dependencies:
RUN  go get -u github.com/AndreasBriese/bbloom \
  && go get -u github.com/dgryski/go-farm \ 
  && go get -u github.com/pkg/errors \
  && go get -u golang.org/x/net/trace \
  && go get -u golang.org/x/sys/unix \
  && go get -u github.com/golang/protobuf/proto

WORKDIR /go/src/github.com/dgraph-io
RUN  git config --global advice.detachedHead false \
  && git clone --branch v1.5.4 --single-branch https://github.com/dgraph-io/badger \
  && cd badger \
  && go build .
# end of "badger"
# - - - - - -

# Get to the tag of go-ipfs, and build all its binaries (gx)
WORKDIR /go/src/github.com/ipfs
RUN git clone -b 'v0.4.19' --single-branch --depth=1 https://github.com/ipfs/go-ipfs \
    && cd go-ipfs \
    && make install
# verify that IPFS can be built
WORKDIR /go/src/github.com/ipfs/go-ipfs/cmd/ipfswatch
RUN go build .

# Copy surgical extraction tool (pre-built in another container)
COPY --from=surgical /go/bin/surgical-extraction /go/bin/surgical-extraction

WORKDIR /go/src/github.com/AtlantPlatform
# Clone fork of go-ipfs
RUN git clone --depth=1 https://github.com/AtlantPlatform/go-ipfs
WORKDIR /go/src/github.com/AtlantPlatform/go-ipfs
# Overwrite makefile and patches files
RUN rm -rf patches/* Makefile
COPY Makefile  ./
COPY patches/*.patch ./patches/

# Cleaning all files in the repository, we can even reset GIT index
RUN make clean
RUN make extract-apply

# `quic-go` uses old structure reference from `mint`
# mint uses QTLS which is 

RUN    go get -u golang.org/x/crypto/chacha20poly1305 \
    && go get -u golang.org/x/crypto/cryptobyte \
    && go get -u golang.org/x/crypto/curve25519 \
    && go get -u golang.org/x/crypto/hkdf

WORKDIR /go/src/github.com/AtlantPlatform/go-ipfs/vendor/github.com/marten-seemann
RUN git clone https://github.com/marten-seemann/qtls \
      && cd ./qtls \
      && git checkout 3637b8f41961577d41c232ecd31ea6791438aeed \
      && go build .

WORKDIR /go/src/github.com/AtlantPlatform/go-ipfs/vendor/github.com/bifurcation
RUN git clone https://github.com/bifurcation/mint \
    && cd ./mint \
    && git checkout 824af65410658916142a7600349144e1289f2110 \
    && go build .

RUN    go get -u github.com/cheekybits/genny/generic \
    && go get -u github.com/hashicorp/golang-lru \
    && go get -u github.com/lucas-clemente/aes12 \
    && go get -u github.com/lucas-clemente/quic-go-certificates

WORKDIR /go/src/github.com/AtlantPlatform/go-ipfs/vendor
RUN find . -type d -name .git | xargs rm -rf

WORKDIR /go/src/github.com/AtlantPlatform/go-ipfs
# RUN make test

# volume for source exchange
VOLUME /output
CMD ["/bin/bash"]