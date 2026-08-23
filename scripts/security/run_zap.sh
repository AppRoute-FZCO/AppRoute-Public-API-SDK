#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)"
TARGET="${ZAP_TARGET_URL:-}"
ISOLATED="${ZAP_ISOLATED:-false}"
ENVIRONMENT="${ZAP_ENVIRONMENT:-}"
OUT_DIR="${ZAP_OUTPUT_DIR:-$ROOT/.security/zap}"

if [[ "$ISOLATED" != "true" || "$ENVIRONMENT" != "sandbox" ]]; then
  echo "Refusing ZAP: set ZAP_ISOLATED=true and ZAP_ENVIRONMENT=sandbox" >&2
  exit 2
fi
if [[ -z "$TARGET" ]]; then
  echo "ZAP_TARGET_URL is required" >&2
  exit 2
fi
if [[ "$TARGET" =~ [Pp][Rr][Oo][Dd] || "$TARGET" =~ [Pp][Rr][Oo][Dd][Uu][Cc][Tt][Ii][Oo][Nn] ]]; then
  echo "Refusing ZAP target that looks like production: $TARGET" >&2
  exit 2
fi
if [[ ! "$TARGET" =~ ^https?://(127\.0\.0\.1|localhost|\[::1\]|[^/]*\.sandbox\.) ]]; then
  echo "Refusing non-local/non-sandbox ZAP target: $TARGET" >&2
  exit 2
fi

mkdir -p "$OUT_DIR"
if command -v zap-baseline.py >/dev/null 2>&1; then
  zap-baseline.py -t "$TARGET" -m 5 -r "$OUT_DIR/zap.html" -J "$OUT_DIR/zap.json"
elif command -v docker >/dev/null 2>&1; then
  docker run --rm --network host -v "$OUT_DIR:/zap/wrk:rw" \
    zaproxy/zap-stable:2.16.1 zap-baseline.py \
    -t "$TARGET" -m 5 -r /zap/wrk/zap.html -J /zap/wrk/zap.json
else
  echo "OWASP ZAP is required (install it or provide Docker)" >&2
  exit 127
fi
