#!/usr/bin/env bash
# Extracts the generated OpenAPI 3.0 spec (JSON + YAML) from a published
# Sapasora Docker image, without needing to run/build the app locally.
#
# Usage:
#   ./scripts/extract-openapi.sh [image-tag] [output-dir]
#
# Examples:
#   ./scripts/extract-openapi.sh                       # latest -> ./openapi-export
#   ./scripts/extract-openapi.sh v2026.40.0             # pin a version
#   ./scripts/extract-openapi.sh latest ./docs/openapi  # custom output dir

set -euo pipefail

IMAGE_TAG="${1:-latest}"
OUTPUT_DIR="${2:-./openapi-export}"
IMAGE="ghcr.io/mrrizkin/sapasora:${IMAGE_TAG}"

echo "Pulling ${IMAGE}..."
docker pull "${IMAGE}"

mkdir -p "${OUTPUT_DIR}"

echo "Extracting OpenAPI spec to ${OUTPUT_DIR}..."
docker run --rm \
  -v "$(cd "${OUTPUT_DIR}" && pwd):/out" \
  --entrypoint sh \
  "${IMAGE}" \
  -c "cp -r /app/public/docs/v3/* /out/"

echo "Done:"
ls -la "${OUTPUT_DIR}"
