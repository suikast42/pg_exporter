#!/bin/bash
#==============================================================#
# File      :   docker/release.sh
# Desc      :   Build and release multi-arch Docker images
# Mtime     :   2025-07-17
# License   :   Apache-2.0 @ https://github.com/pgsty/pg_exporter
# Copyright :   2018-2026  Ruohang Feng / Vonng (rh@vonng.com)
#==============================================================#

set -euo pipefail

# Get current script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "${SCRIPT_DIR}")"

cd "${PROJECT_ROOT}"

# Get version from Makefile or environment
VERSION=${VERSION:-$(grep '^VERSION' Makefile | cut -d'=' -f2 | tr -d ' ?')}
DOCKER_REPO=${DOCKER_REPO:-pgsty/pg_exporter}
PLATFORMS=${PLATFORMS:-linux/amd64,linux/arm64}
PUSH=${PUSH:-true}

echo "Building and releasing multi-arch Docker images for pg_exporter ${VERSION}"
echo "Platforms: ${PLATFORMS}"
echo "Repository: ${DOCKER_REPO}"

# Check if both arch binaries exist
echo "Checking if release binaries exist..."
if [[ ! -d "dist/${VERSION}" ]]; then
    echo "Error: dist/${VERSION} directory not found. Building Linux binaries first..."
    make release-linux
fi

AMD64_TAR="dist/${VERSION}/pg_exporter-${VERSION}.linux-amd64.tar.gz"
ARM64_TAR="dist/${VERSION}/pg_exporter-${VERSION}.linux-arm64.tar.gz"

if [[ ! -f "${AMD64_TAR}" ]] || [[ ! -f "${ARM64_TAR}" ]]; then
    echo "Error: Required binaries not found. Building them now..."
    make release-linux
    if [[ ! -f "${AMD64_TAR}" ]] || [[ ! -f "${ARM64_TAR}" ]]; then
        echo "Error: Failed to build required binaries"
        exit 1
    fi
fi

# Create temporary build contexts for each architecture
BUILD_DIR=$(mktemp -d)
trap "rm -rf ${BUILD_DIR}" EXIT

echo "Created temporary build directory: ${BUILD_DIR}"

# Extract binaries for both architectures
echo "Extracting binaries..."

mkdir -p "${BUILD_DIR}/amd64" "${BUILD_DIR}/arm64"

tar -xzf "${AMD64_TAR}" -C "${BUILD_DIR}/amd64" --strip-components=1
tar -xzf "${ARM64_TAR}" -C "${BUILD_DIR}/arm64" --strip-components=1

# Create multi-arch Dockerfile that works with buildx
cat > "${BUILD_DIR}/Dockerfile" << 'EOF'
FROM scratch

LABEL org.opencontainers.image.authors="Ruohang Feng <rh@vonng.com>" \
      org.opencontainers.image.url="https://github.com/pgsty/pg_exporter" \
      org.opencontainers.image.source="https://github.com/pgsty/pg_exporter" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.title="pg_exporter" \
      org.opencontainers.image.description="PostgreSQL/Pgbouncer metrics exporter for Prometheus"

WORKDIR /bin
COPY pg_exporter /bin/pg_exporter
COPY pg_exporter.yml /etc/pg_exporter.yml
COPY LICENSE /LICENSE

EXPOSE 9630/tcp
ENTRYPOINT ["/bin/pg_exporter"]
EOF

# Copy Dockerfile to both arch directories
cp "${BUILD_DIR}/Dockerfile" "${BUILD_DIR}/amd64/Dockerfile"
cp "${BUILD_DIR}/Dockerfile" "${BUILD_DIR}/arm64/Dockerfile"

echo "Setting up Docker buildx..."

# Create or use existing buildx builder
BUILDER_NAME="pg_exporter_builder"
if ! docker buildx ls | grep -q "${BUILDER_NAME}"; then
    echo "Creating new buildx builder: ${BUILDER_NAME}"
    docker buildx create --name "${BUILDER_NAME}" --use --bootstrap
else
    echo "Using existing buildx builder: ${BUILDER_NAME}"
    docker buildx use "${BUILDER_NAME}"
fi

if [[ "${PUSH}" == "true" ]]; then
    echo "Building and pushing multi-arch images..."
    echo "This will create a manifest list that automatically selects the right architecture."
    
    # Build and push multi-arch images with manifest list
    # This creates the "fat manifest" that allows automatic architecture selection
    docker buildx build \
        --platform "${PLATFORMS}" \
        --file "${BUILD_DIR}/amd64/Dockerfile" \
        --tag "${DOCKER_REPO}:${VERSION}" \
        --tag "${DOCKER_REPO}:latest" \
        --push \
        "${BUILD_DIR}/amd64"
else
    echo "Building multi-arch images locally (no push)..."
    echo "This will create local images for testing."
    
    # Build multi-arch images locally without pushing
    docker buildx build \
        --platform "${PLATFORMS}" \
        --file "${BUILD_DIR}/amd64/Dockerfile" \
        --tag "${DOCKER_REPO}:${VERSION}" \
        --tag "${DOCKER_REPO}:latest" \
        --load \
        "${BUILD_DIR}/amd64"
fi

echo ""
if [[ "${PUSH}" == "true" ]]; then
    echo "✅ Multi-arch Docker images released successfully!"
    echo ""
    echo "Images pushed:"
    echo "  ${DOCKER_REPO}:${VERSION}"
    echo "  ${DOCKER_REPO}:latest"
else
    echo "✅ Multi-arch Docker images built locally!"
    echo ""
    echo "Images created:"
    echo "  ${DOCKER_REPO}:${VERSION}"
    echo "  ${DOCKER_REPO}:latest"
fi
echo ""
echo "These images support the following architectures:"
echo "  - linux/amd64"
echo "  - linux/arm64"
echo ""
if [[ "${PUSH}" == "true" ]]; then
    echo "Users can now pull without specifying architecture:"
    echo "  docker pull ${DOCKER_REPO}:${VERSION}"
    echo "  docker pull ${DOCKER_REPO}:latest"
    echo ""
    echo "Docker will automatically select the correct architecture for their platform."
    
    # Verify the manifest
    echo ""
    echo "Verifying manifest list..."
    docker buildx imagetools inspect "${DOCKER_REPO}:${VERSION}" || echo "Note: manifest inspection requires authentication"
else
    echo "To test the local images:"
    echo "  docker run --rm ${DOCKER_REPO}:${VERSION} --version"
    echo ""
    echo "To push later:"
    echo "  docker push ${DOCKER_REPO}:${VERSION}"
    echo "  docker push ${DOCKER_REPO}:latest"
fi