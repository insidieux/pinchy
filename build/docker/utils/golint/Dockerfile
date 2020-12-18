ARG GO_VERSION=1.15
FROM golang:${GO_VERSION}

RUN go get -u -t golang.org/x/lint/golint

ENTRYPOINT ["/go/bin/golint"]
