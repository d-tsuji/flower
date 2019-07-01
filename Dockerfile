# Build Stage
FROM golang:1.12.5 AS builder
ENV REPOSITORY github.com/d-tsuji/flower
ADD . $GOPATH/src/$REPOSITORY
WORKDIR $GOPATH/src/$REPOSITORY
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-s -w' -a -installsuffix cgo -o /main main.go

# Optional
ARG DIR=$GOPATH/src/$REPOSITORY
COPY ./wait-for-postgres.sh $DIR/wait-for-postgres.sh
RUN chmod +x $DIR/wait-for-postgres.sh

# Runtime Stage
FROM alpine:3.9.4
RUN apk add --no-cache ca-certificates
COPY --from=builder /main .
# CMD ["./main", "--host", "0.0.0.0"]
CMD ["./wait-for-postgres.sh", "./main", "--host", "0.0.0.0"]

EXPOSE 8021