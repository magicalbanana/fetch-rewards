FROM golang:1.13.0-alpine3.10 as builder

RUN mkdir -p /go/src/github.com/magicalbanana/fetch-rewards/

WORKDIR /go/src/github.com/magicalbanana/fetch-rewards/

RUN apk add --update --no-cache alpine-sdk git

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -mod vendor -v -a -installsuffix cgo -o http-server \
    cmd/fetch-rewards-http-server/fetch-rewards-http-server.go

# actual container
FROM alpine:3.10

RUN apk add --update --no-cache bash ca-certificates

RUN mkdir -p /app

WORKDIR /app


COPY resources/jsonschema/ ./resources/jsonschema
COPY resources/postgres/sqlstatements/ ./resources/postgres/sqlstatements

RUN mkdir -p scripts/

COPY scripts/run-http-server.sh ./scripts/

COPY --from=builder /go/src/github.com/magicalbanana/fetch-rewards/http-server .

CMD ["./scripts/docker/run-http-server.sh"]
