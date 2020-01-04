############################################################
# Step.1 Build Stage
############################################################
FROM golang:1.13 AS builder
ENV REPOSITORY github.com/d-tsuji/flower-v2
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
ADD . $GOPATH/src/$REPOSITORY
WORKDIR $GOPATH/src/$REPOSITORY
RUN go build -ldflags '-s -w' -a -installsuffix cgo -o /register cmd/register/main.go

############################################################
# Step.2 Runtime Stage
############################################################
FROM alpine:3.11.2
RUN apk add --no-cache ca-certificates

## dockerize(waiting for running postgreSql)
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

COPY --from=builder /register .