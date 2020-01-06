############################################################
# Step.1 Build Stage
############################################################
FROM golang:1.13 AS builder
ARG APP_NAME
ENV REPOSITORY github.com/d-tsuji/flower
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
WORKDIR $GOPATH/src/$REPOSITORY
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags '-s -w' -a -installsuffix cgo -o /$APP_NAME cmd/$APP_NAME/main.go

############################################################
# Step.2 Runtime Stage
############################################################
FROM alpine:3.11.2
ARG APP_NAME
RUN apk add --no-cache ca-certificates

## dockerize(waiting for running postgreSql)
ENV DOCKERIZE_VERSION v0.6.1
#RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
#    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
#    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz
COPY assets/tools/dockerize-linux-amd64-$DOCKERIZE_VERSION/dockerize /usr/local/bin

COPY --from=builder /$APP_NAME .