# Build surgical in a separate docker image
FROM golang:1.11.5-stretch as surgical
WORKDIR /go/src/github.com/AtlantPlatform
COPY ./ ./surgical-extraction/

RUN  go get -u "github.com/jawher/mow.cli" \
  && go get -u "github.com/oklog/ulid" \
  && go get -u "github.com/xlab/closer"

WORKDIR /go/src/github.com/AtlantPlatform/surgical-extraction
RUN go build -a -installsuffix cgo --ldflags "-s -w" . && go install
CMD ["which", "surgical-extraction"]