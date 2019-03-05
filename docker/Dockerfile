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

# - - - - - - - - - - - - - - - - - - 
# rebuild all QUIC transport from the source
# its dependencies seems to be broken/lost compatibility with master
# ENV GO111MODULE=on
# WORKDIR $GOPATH/src/github.com/marten-seemann
# RUN git clone https://github.com/marten-seemann/qtls \
#      && cd ./qtls \
#      && git checkout 3637b8f41961577d41c232ecd31ea6791438aeed \
#      && go get -v ./... \
#      && go build .

# WORKDIR $GOPATH/src/github.com/bifurcation
# RUN git clone https://github.com/bifurcation/mint \
#     && cd ./mint \
#     && go get -v ./... \
#     && go build .

# WORKDIR $GOPATH/src/github.com/lucas-clemente
# RUN git clone https://github.com/lucas-clemente/quic-go \
#     && cd ./quic-go \
#     && go get -v ./... \
#     && go build .

# WORKDIR $GOPATH/src/github.com/libp2p \
#     && git clone https://github.com/libp2p/go-libp2p-quic-transport \
#     && cd ./go-libp2p-quic-transport \
#     && go get -v ./... \
#     && go build .

# end of "quic"
# - - - - - - - - - - - - - - - - - - 

# Get fresh go-ipfs, and build all its binaries (gx)
ENV GO111MODULE=off
RUN go get -u github.com/ipfs/go-ipfs
WORKDIR /go/src/github.com/ipfs/go-ipfs
RUN make install
# verify that IPFS can be build
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
# Replace some local packages with github-fetched (masters or branches)
# RUN make replace
# RUN make test

# volume for source exchange
VOLUME /output
CMD ["/bin/bash"]