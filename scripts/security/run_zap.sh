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
# Parse the URL instead of applying a substring regex. This rejects userinfo
# tricks (trusted.sandbox.example@evil.example), look-alike suffixes
# (foo.sandbox.example.evil.com), and malformed local-host prefixes.
python3 - "$TARGET" <<'PYURL'
import re
import sys
from urllib.parse import urlsplit

raw = sys.argv[1]
try:
    parsed = urlsplit(raw)
    host = parsed.hostname
    _port = parsed.port
except ValueError:
    parsed = None
    host = None

if parsed is None or parsed.scheme not in {"http", "https"} or not host:
    raise SystemExit(2)
if parsed.username is not None or parsed.password is not None:
    raise SystemExit(2)
if parsed.fragment:
    raise SystemExit(2)

host = host.rstrip(".").lower()
if host not in {"localhost", "127.0.0.1", "::1"}:
    labels = host.split(".")
    # Explicit three-label sandbox form: <service>.sandbox.<domain>.
    if len(labels) != 3 or labels[1] != "sandbox" or any(
        not re.fullmatch(r"[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?", label)
        for label in labels
    ):
        raise SystemExit(2)
PYURL

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
