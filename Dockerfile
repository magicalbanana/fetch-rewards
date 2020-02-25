FROM golang:1.13.8-alpine3.11 as builder

RUN mkdir -p /go/src/github.com/magicalbanana/fetch-rewards/

WORKDIR /go/src/github.com/magicalbanana/fetch-rewards/

COPY . .

RUN apk add --update --no-cache alpine-sdk git

RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -mod vendor -v -a -installsuffix cgo -o http-server \
    main.go

# actual container
FROM alpine:3.11

RUN apk add --update --no-cache bash ca-certificates

RUN mkdir -p /app

WORKDIR /app

RUN mkdir -p scripts/

COPY scripts/run-http-server.sh ./scripts/

COPY --from=builder /go/src/github.com/magicalbanana/fetch-rewards/http-server .

CMD ["./scripts/run-http-server.sh"]
