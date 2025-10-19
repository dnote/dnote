#!/usr/bin/env bash
# Build Docker image for local testing
#
# Usage:
# # builds for host platform (auto-detected)
# ./build.sh 1.0.0
#
# # builds arm64
# ./build.sh 1.0.0 linux/arm64
#
# # builds multiple platforms
# ./build.sh 1.0.0 "linux/amd64,linux/arm64,linux/arm/v7,linux/386"
set -eux

version=$1

# Detect host platform if not specified
if [ -z "${2:-}" ]; then
  HOST_ARCH=$(uname -m)
  case "$HOST_ARCH" in
    x86_64)        platform="linux/amd64" ;;
    aarch64|arm64) platform="linux/arm64" ;;
    armv7l)        platform="linux/arm/v7" ;;
    i386|i686)     platform="linux/386" ;;
    *)
      echo "Warning: Unsupported architecture: $HOST_ARCH, defaulting to linux/amd64"
      platform="linux/amd64"
      ;;
  esac
  echo "Auto-detected platform: $platform"
else
  platform=$2
fi

dir=$(dirname "${BASH_SOURCE[0]}")
projectDir="$dir/../.."

# Copy all Linux tarballs to Docker build context
cp "$projectDir/build/server/dnote_server_${version}_linux_amd64.tar.gz" "$dir/"
cp "$projectDir/build/server/dnote_server_${version}_linux_arm64.tar.gz" "$dir/"
cp "$projectDir/build/server/dnote_server_${version}_linux_arm.tar.gz" "$dir/"
cp "$projectDir/build/server/dnote_server_${version}_linux_386.tar.gz" "$dir/"

# Count platforms (check for comma)
if [[ "$platform" == *","* ]]; then
  echo "Building for multiple platforms: $platform"

  # Check if multiarch builder exists, create if not
  if ! docker buildx ls | grep -q "multiarch"; then
    echo "Creating multiarch builder for multi-platform builds..."
    docker buildx create --name multiarch --use
    docker buildx inspect --bootstrap
  else
    echo "Using existing multiarch builder"
    docker buildx use multiarch
  fi
  echo ""

  docker buildx build \
    --platform "$platform" \
    -t dnote/dnote:"$version" \
    -t dnote/dnote:latest \
    --build-arg version="$version" \
    "$dir"

  # Switch back to default builder
  docker buildx use default
else
  echo "Building for single platform: $platform"
  echo "Image will be loaded to local Docker daemon"

  docker buildx build \
    --platform "$platform" \
    -t dnote/dnote:"$version" \
    -t dnote/dnote:latest \
    --build-arg version="$version" \
    --load \
    "$dir"
fi

echo ""
echo "Build complete!"
if [[ "$platform" != *","* ]]; then
  echo "Test with: docker run --rm dnote/dnote:$version ./dnote-server version"
fi
