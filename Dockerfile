# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.25.1
ARG NODE_VERSION=22-bookworm-slim
ARG TDLIB_COMMIT=971684a3dcc7bdf99eec024e1c4f57ae729d6d53

FROM node:${NODE_VERSION} AS assets

WORKDIR /src

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN corepack enable && corepack prepare pnpm@12.4.2 --activate \
    && pnpm install --frozen-lockfile

COPY . .
RUN pnpm build:assets

FROM golang:${GO_VERSION}-bookworm AS builder

ARG TDLIB_COMMIT
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        build-essential \
        clang \
        cmake \
        gperf \
        git \
        libc++-dev \
        libc++abi-dev \
        libssl-dev \
        php-cli \
        zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*

# go-tdlib requires the matching native TDLib version.
RUN git clone --depth=1 https://github.com/tdlib/td.git /tmp/td \
    && cd /tmp/td \
    && git fetch --depth=1 origin "${TDLIB_COMMIT}" \
    && git checkout "${TDLIB_COMMIT}" \
    && cmake -S . -B build \
        -DCMAKE_BUILD_TYPE=Release \
        -DCMAKE_INSTALL_PREFIX=/usr/local \
    && cmake --build build --target install -j"$(nproc)" \
    && rm -rf /tmp/td

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=assets /src/public/build ./public/build

RUN go tool templ generate
RUN CGO_ENABLED=1 go build -tags libtdjson -trimpath -ldflags="-s -w" -o /out/app ./cmd/main
RUN CGO_ENABLED=1 go build -tags libtdjson -trimpath -ldflags="-s -w" -o /out/toolbox ./cmd/toolbox

FROM debian:bookworm-slim AS runtime

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        libgcc-s1 \
        libssl3 \
        libstdc++6 \
        zlib1g \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /usr/local/lib/libtdjson.so* /usr/local/lib/
COPY --from=builder /out/app /usr/local/bin/sapasora
COPY --from=builder /out/toolbox /usr/local/bin/toolbox
COPY --from=builder /src/public ./public

RUN mkdir -p /app/storage /app/.tdlib

ENV LD_LIBRARY_PATH=/usr/local/lib
ENV APP_PORT=3000

EXPOSE 3000

CMD ["/usr/local/bin/sapasora"]
