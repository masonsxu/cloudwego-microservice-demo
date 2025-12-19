#!/bin/bash

# Exit on error
set -e

# Get the directory of the script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# The project root is three levels up from the script's directory
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

# Set the default image name
DEFAULT_IMAGE_NAME="gateway"

# Use environment variables or defaults
IMAGE_NAME="${IMAGE_NAME:-$DEFAULT_IMAGE_NAME}"
CONTAINER_CMD="${CONTAINER_CMD:-docker}"

# Generate a timestamp for the image tag
TIMESTAMP_TAG=$(date +%Y%m%d-%H%M%S)

# The Dockerfile path, relative to the project root
DOCKERFILE="gateway/docker/Dockerfile"

# Define the full image names
TIMESTAMP_IMAGE="$IMAGE_NAME:$TIMESTAMP_TAG"
LATEST_IMAGE="$IMAGE_NAME:latest"

echo "Using container command: $CONTAINER_CMD"
echo "Building container image: $TIMESTAMP_IMAGE"
echo "Project root (build context): $PROJECT_ROOT"
echo "Dockerfile: $PROJECT_ROOT/$DOCKERFILE"

# Change to the project root directory to execute the build
cd "$PROJECT_ROOT"

# Build the container image with the timestamp tag
if [ "$CONTAINER_CMD" = "podman" ]; then
  "$CONTAINER_CMD" build --format docker -t "$TIMESTAMP_IMAGE" -f "$DOCKERFILE" .
else
  "$CONTAINER_CMD" build -t "$TIMESTAMP_IMAGE" -f "$DOCKERFILE" .
fi

# Tag the new image as 'latest'
"$CONTAINER_CMD" tag "$TIMESTAMP_IMAGE" "$LATEST_IMAGE"

echo "Container image built and tagged successfully:"
echo "  - $TIMESTAMP_IMAGE"
echo "  - $LATEST_IMAGE"
