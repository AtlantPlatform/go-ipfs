# Build surgical in a separate docker image
FROM golang:1.12-stretch as surgical
WORKDIR /go/src/github.com/AtlantPlatform/surgical-extraction
COPY ./vendor/ ./vendor/
COPY go.mod .
COPY go.sum .
COPY *.go ./

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor 
RUN go build && go install
CMD ["which", "surgical-extraction"]