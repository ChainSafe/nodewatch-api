FROM  golang:1.16-alpine3.13 AS builder

WORKDIR /code

COPY go.mod .
COPY go.sum .
RUN go mod download

# build the binary
ADD . .
RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /crawler cmd/main.go

# final stage
FROM alpine:3.14.0

ARG env=dev

RUN apk add curl
COPY --from=builder /crawler /
COPY cmd/config/config.$env.yaml /

RUN chmod +x /crawler
