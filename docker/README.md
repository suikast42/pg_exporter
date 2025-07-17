# Docker Build Scripts

This directory contains scripts for building Docker images for pg_exporter.

## Scripts

### `build.sh` - Local Development Build

Builds a single-architecture Docker image for local development and testing.

```bash
# Build for current architecture
./docker/build.sh

# Build for specific architecture
ARCH=arm64 ./docker/build.sh

# Use custom repository
DOCKER_REPO=myrepo/pg_exporter ./docker/build.sh
```

**Features:**
- Detects current architecture automatically
- Builds only for the current platform
- Creates local tags: `<repo>:dev`, `<repo>:latest-<arch>`, `<repo>:<version>-<arch>`
- Does not push to registry
- Automatically builds missing Linux binaries if needed

### `release.sh` - Production Multi-Arch Release

Builds and pushes multi-architecture Docker images with manifest list support.

```bash
# Build and push multi-arch images
./docker/release.sh

# Build locally without pushing (for testing)
PUSH=false ./docker/release.sh

# Custom repository
DOCKER_REPO=pgsty/pg_exporter ./docker/release.sh

# Custom platforms
PLATFORMS=linux/amd64,linux/arm64,linux/arm/v7 ./docker/release.sh
```

**Features:**
- Builds for multiple architectures (amd64, arm64 by default)
- Creates manifest list for automatic architecture selection
- Pushes to Docker registry
- Creates tags: `<repo>:<version>`, `<repo>:latest`
- Requires pre-built Linux binaries (`make release-linux`)

## How Multi-Arch Works

The `release.sh` script uses Docker buildx to create a **manifest list** (also called "fat manifest"). This allows users to pull images without specifying architecture:

```bash
# Users can simply run:
docker pull pgsty/pg_exporter:latest

# Docker automatically selects the right architecture:
# - On AMD64 systems: pulls linux/amd64 image
# - On ARM64 systems: pulls linux/arm64 image
```

## Prerequisites

### For Local Build (`build.sh`)
- Docker
- Make (for building binaries if needed)

### For Production Release (`release.sh`)
- Docker with buildx support
- Docker registry authentication (for pushing)
- Pre-built Linux binaries: `make release-linux`

## Environment Variables

| Variable      | Default                   | Description                              |
|---------------|---------------------------|------------------------------------------|
| `VERSION`     | From Makefile             | Image version tag                        |
| `DOCKER_REPO` | `pgsty/pg_exporter`       | Docker repository                        |
| `ARCH`        | Auto-detected             | Target architecture (build.sh only)      |
| `PLATFORMS`   | `linux/amd64,linux/arm64` | Target platforms (release.sh only)       |
| `PUSH`        | `true`                    | Whether to push images (release.sh only) |

## Examples

```bash
# Local development
./docker/build.sh
docker run --rm pgsty/pg_exporter:dev --version

# Production release
make release-linux  # Build binaries first
./docker/release.sh

# Test locally without pushing
PUSH=false ./docker/release.sh

# Custom repository and platforms
DOCKER_REPO=mycompany/pg_exporter \
PLATFORMS=linux/amd64,linux/arm64,linux/arm/v7 \
./docker/release.sh
```