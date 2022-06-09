FROM alpine:3.15.0
RUN apk add --update --no-cache git
ENTRYPOINT ["/go/bin/go-bpkg"]
COPY go-bpkg /go/bin/go-bpkg
