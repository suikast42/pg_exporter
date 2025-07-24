#!/bin/bash
#==============================================================#
# File      :   docker/build.sh
# Desc      :   Build single-arch Docker image locally with Go
# Mtime     :   2025-07-17
# License   :   Apache-2.0 @ https://github.com/pgsty/pg_exporter
# Copyright :   2018-2025  Ruohang Feng / Vonng (rh@vonng.com)
#==============================================================#

set -euo pipefail

# Get current script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "${SCRIPT_DIR}")"

cd "${PROJECT_ROOT}"

# Get version from Makefile or environment
VERSION=${VERSION:-$(grep '^VERSION' Makefile | cut -d'=' -f2 | tr -d ' ?')}
DOCKER_REPO=${DOCKER_REPO:-vonng/pg_exporter}

echo "Building Docker image for pg_exporter ${VERSION} - LOCAL BUILD (Go source)"

# Create Dockerfile for local Go build
cat > "${SCRIPT_DIR}/Dockerfile.local" << 'EOF'
# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder-env

WORKDIR /build

# Copy dependency files and download deps
COPY go.mod go.sum ./
RUN \
  --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=linux go mod download

# Copy source code
COPY . /build

# Build static binary
RUN \
  --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /pg_exporter .

FROM scratch

LABEL org.opencontainers.image.authors="Ruohang Feng <rh@vonng.com>" \
      org.opencontainers.image.url="https://github.com/pgsty/pg_exporter" \
      org.opencontainers.image.source="https://github.com/pgsty/pg_exporter" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.title="pg_exporter" \
      org.opencontainers.image.description="PostgreSQL/Pgbouncer metrics exporter for Prometheus"

WORKDIR /bin
COPY --from=builder-env /pg_exporter /bin/pg_exporter
COPY pg_exporter.yml /etc/pg_exporter.yml

EXPOSE 9630/tcp
ENTRYPOINT ["/bin/pg_exporter"]
EOF

echo "Building Docker image with Go source..."

# Build image using Go source (local only, no push)
docker build \
    -f "${SCRIPT_DIR}/Dockerfile.local" \
    -t "${DOCKER_REPO}:${VERSION}-dev" \
    -t "${DOCKER_REPO}:dev" \
    "${PROJECT_ROOT}"

echo "Docker image built successfully:"
echo "  ${DOCKER_REPO}:${VERSION}-dev"
echo "  ${DOCKER_REPO}:dev"
echo ""
echo "To test the image:"
echo "  docker run --rm ${DOCKER_REPO}:dev --version"

# Clean up
rm -f "${SCRIPT_DIR}/Dockerfile.local"