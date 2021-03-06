FROM golang:1.18-alpine AS build

RUN apk add --update git
RUN apk add ca-certificates

WORKDIR /go/src/github.com/rafaelcalleja/go-bpkg

COPY . .

RUN go mod tidy && TAG=$(git describe --tags --abbrev=0) \
    && LDFLAGS=$(echo "-s -w -X main.version="$TAG) \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/go-bpkg -ldflags "$LDFLAGS" cmd/main.go

# Building image with the binary
FROM scratch

COPY --from=build /go/bin/go-bpkg /go/bin/go-bpkg
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/go/bin/go-bpkg"]
