# Build Stage
FROM golang:1.13 AS builder
ENV REPOSITORY github.com/d-tsuji/flower-v2
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ADD ./wait-for-postgres.sh /opt/
ADD . $GOPATH/src/$REPOSITORY
WORKDIR $GOPATH/src/$REPOSITORY
RUN go build -ldflags '-s -w' -a -installsuffix cgo -o /watcher cmd/watcher/main.go

# Runtime Stage
FROM alpine:3.11.2
RUN apk add --no-cache ca-certificates
COPY --from=builder /watcher .
CMD ["./watcher", "--host", "0.0.0.0"]