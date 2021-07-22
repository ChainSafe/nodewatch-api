# Copyright 2021 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

FROM golang:1.16-alpine AS builder

RUN apk add build-base
WORKDIR /code
COPY go.mod .
COPY go.sum .
RUN go mod download

# build the binary
ADD . .
RUN env GOOS=linux GOARCH=amd64 go build -o /crawler cmd/main.go

# final stage
FROM alpine:3.14.0

RUN apk add build-base
ARG env=dev

RUN apk add curl
COPY --from=builder /crawler /
COPY cmd/config/config.$env.yaml /config.yaml

RUN chmod +x /crawler
ENTRYPOINT ["/crawler", "-p", "/config.yaml"]
