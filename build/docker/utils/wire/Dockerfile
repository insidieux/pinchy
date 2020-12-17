ARG GO_VERSION=1.15
FROM golang:${GO_VERSION}

RUN go get -u -t github.com/google/wire/cmd/wire

ENTRYPOINT ["/go/bin/wire"]
