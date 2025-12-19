#!/bin/bash

# Exit on error
set -e

# Get the directory of the script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Set the default image name
DEFAULT_IMAGE_NAME="identity-srv"

# Use environment variables or defaults
IMAGE_NAME="${IMAGE_NAME:-$DEFAULT_IMAGE_NAME}"
CONTAINER_CMD="${CONTAINER_CMD:-docker}"

# Generate a timestamp for the image tag
TIMESTAMP_TAG=$(date +%Y%m%d-%H%M%S)

# The build context is the parent directory of the script's location
BUILD_CONTEXT="$DIR/.."

# The Dockerfile is in the same directory as the script
DOCKERFILE="$DIR/Dockerfile"

# Define the full image names
TIMESTAMP_IMAGE="$IMAGE_NAME:$TIMESTAMP_TAG"
LATEST_IMAGE="$IMAGE_NAME:latest"

echo "Using container command: $CONTAINER_CMD"
echo "Building container image: $TIMESTAMP_IMAGE"
echo "Build context: $BUILD_CONTEXT"
echo "Dockerfile: $DOCKERFILE"

# Build the container image with the timestamp tag
if [ "$CONTAINER_CMD" = "podman" ]; then
  "$CONTAINER_CMD" build --format docker -t "$TIMESTAMP_IMAGE" -f "$DOCKERFILE" "$BUILD_CONTEXT"
else
  "$CONTAINER_CMD" build -t "$TIMESTAMP_IMAGE" -f "$DOCKERFILE" "$BUILD_CONTEXT"
fi

# Tag the new image as 'latest'
"$CONTAINER_CMD" tag "$TIMESTAMP_IMAGE" "$LATEST_IMAGE"

echo "Container image built and tagged successfully:"
echo "  - $TIMESTAMP_IMAGE"
echo "  - $LATEST_IMAGE"
