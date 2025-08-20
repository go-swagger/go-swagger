#!/usr/bin/env bash
set -euo pipefail

SPEC_PATH="${1:-OpenAPI31CodescanMVPSpecification.yml}"

# Validate using Redocly CLI (supports OpenAPI 3.1)
docker run --rm -v "$(pwd)":/work -w /work redocly/cli:latest lint "$SPEC_PATH"

# If you also want to validate the JSON variant, uncomment:
# docker run --rm -v "$(pwd)":/work -w /work redocly/cli:latest lint OpenAPI31CodescanMVPSpecification.json

echo "Validation passed for: $SPEC_PATH"
