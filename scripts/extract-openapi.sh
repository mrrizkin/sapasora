#!/usr/bin/env bash
# Extracts the generated OpenAPI 3.0 spec (JSON + YAML) from a published
# Sapasora Docker image, without needing to run/build the app locally.
# Optionally rewrites the `servers` field to a real domain, since the spec
# is generated with a hardcoded localhost:3000 placeholder by default.
#
# Usage:
#   ./scripts/extract-openapi.sh [image-tag] [output-dir] [server-url]
#
# Examples:
#   ./scripts/extract-openapi.sh
#     -> pulls :latest, writes to ./openapi-export, keeps localhost servers
#
#   ./scripts/extract-openapi.sh latest ./openapi-export https://sapasora.mrrizkin.com
#     -> rewrites the servers field to https://sapasora.mrrizkin.com
#
#   ./scripts/extract-openapi.sh v2026.40.0 ./openapi-export https://sapasora.mrrizkin.com
#     -> same, but pinned to a specific image version

set -euo pipefail

IMAGE_TAG="${1:-latest}"
OUTPUT_DIR="${2:-./openapi-export}"
SERVER_URL="${3:-}"
IMAGE="ghcr.io/mrrizkin/sapasora:${IMAGE_TAG}"

echo "Pulling ${IMAGE}..."
docker pull "${IMAGE}"

mkdir -p "${OUTPUT_DIR}"

echo "Extracting OpenAPI spec to ${OUTPUT_DIR}..."
docker run --rm \
  --user "$(id -u):$(id -g)" \
  -v "$(cd "${OUTPUT_DIR}" && pwd):/out" \
  --entrypoint sh \
  "${IMAGE}" \
  -c "cp -r /app/public/docs/v3/* /out/"

if [ -n "${SERVER_URL}" ]; then
  echo "Rewriting servers field to ${SERVER_URL}..."

  if command -v jq >/dev/null 2>&1; then
    jq --arg url "${SERVER_URL}" '.servers = [{"url": $url}]' \
      "${OUTPUT_DIR}/openapi.json" > "${OUTPUT_DIR}/openapi.json.tmp"
    mv "${OUTPUT_DIR}/openapi.json.tmp" "${OUTPUT_DIR}/openapi.json"
  else
    echo "Warning: jq not found, skipping openapi.json rewrite" >&2
  fi

  python3 - "${OUTPUT_DIR}/openapi.yaml" "${SERVER_URL}" <<'PY'
import sys
import yaml

path, url = sys.argv[1], sys.argv[2]
with open(path) as f:
    spec = yaml.safe_load(f)

spec["servers"] = [{"url": url}]

with open(path, "w") as f:
    yaml.safe_dump(spec, f, sort_keys=False)
PY
fi

echo "Done:"
ls -la "${OUTPUT_DIR}"
