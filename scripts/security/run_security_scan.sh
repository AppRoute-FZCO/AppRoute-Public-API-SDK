#!/usr/bin/env bash
# The same digest-pinned, Linux/amd64 scanners run locally and in CI.
set -euo pipefail
ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)"
RESULTS="$ROOT/.security/results"
GITLEAKS='zricethezav/gitleaks@sha256:bf00b5e039f0fad4b32935dc5ec1e358f227ccd097bcb64b971f0331072fe2ae' # 8.28.0
SEMGREP='semgrep/semgrep@sha256:61a7ab31cdab865212ae8ed3ddefc0a51d61ef3b016a63e79848f56241d4586d' # 1.136.0
TRIVY='aquasec/trivy@sha256:adbf5948d99de471725d37d2cfbd982178971d1cf081615284671d1c8ba99fb8' # 0.66.0

case "${1:-}" in
  gitleaks)
    docker run --rm --platform linux/amd64 --network none -v "$ROOT:/repo:ro" "$GITLEAKS" \
      dir /repo --redact=100 --config /repo/.security/gitleaks.toml --no-banner
    ;;
  semgrep)
    docker run --rm --platform linux/amd64 --network none -v "$ROOT:/src:ro" "$SEMGREP" semgrep \
      scan --config /src/.security/semgrep.yml --error --strict --no-git-ignore --metrics=off --disable-version-check \
      --exclude .worktrees --exclude .venv --exclude node_modules --exclude .security/results /src
    ;;
  trivy-config)
    docker run --rm --platform linux/amd64 -v "$ROOT:/repo:ro" "$TRIVY" \
      fs --scanners misconfig --severity HIGH,CRITICAL --exit-code 1 \
      --skip-dirs /repo/.git,/repo/.worktrees,/repo/.venv,/repo/node_modules,/repo/.security/results /repo
    ;;
  trivy-vuln)
    mkdir -p "$RESULTS"
    # The report is recreated for each run. A failed tool or missing report
    # must not be confused with a successful scan containing dependency CVEs.
    rm -f "$RESULTS/trivy-vuln.json"
    docker run --rm --platform linux/amd64 -v "$ROOT:/repo:ro" -v "$RESULTS:/reports" "$TRIVY" \
      fs --scanners vuln --severity HIGH,CRITICAL --exit-code 0 \
      --format json --output /reports/trivy-vuln.json \
      --skip-dirs /repo/.git,/repo/.worktrees,/repo/.venv,/repo/node_modules,/repo/.security/results /repo
    test -s "$RESULTS/trivy-vuln.json"
    ;;
  *)
    echo 'Usage: scripts/security/run_security_scan.sh {gitleaks|semgrep|trivy-config|trivy-vuln}' >&2
    exit 2
    ;;
esac
