
FROM golang:1.23.9-bullseye AS builder
RUN apt-get update && apt-get upgrade -y && apt-get clean
RUN apt-get install -y bluetooth bluez libbluetooth-dev
WORKDIR /workspace

# ビルドに必要なファイルを全てコピー
COPY go.mod go.sum ./   
RUN go mod download

COPY . .         
# 静的リンク (musl libc) にしたいなら CGO_ENABLED=0 もここで指定
RUN go build -o peekbt ./cmd/main

FROM alpine:latest
ARG VERSION=0.2.5
LABEL org.opencontainers.image.source="https://github.com/ozsys/peekbt" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.title="peekbt" \
      org.opencontainers.image.description="CLI tool for Bluetooth scan & info"

# 非 root ユーザーを用意
RUN adduser -D -h /work nonroot
WORKDIR /work

# ビルド産物をコピー
COPY --from=builder /workspace/peekbt /usr/local/bin/peekbt

USER nonroot
ENTRYPOINT ["peekbt"]
