#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)"; cd "$ROOT"
have(){ command -v "$1" >/dev/null 2>&1; }
run_gitleaks(){ if have gitleaks; then gitleaks detect --source . --redact --config .security/gitleaks.toml; else docker run --rm -v "$ROOT:/repo:ro" zricethezav/gitleaks:v8.28.0 detect --source=/repo --redact --config=/repo/.security/gitleaks.toml --no-banner; fi; }
run_semgrep(){ if have semgrep; then semgrep scan --config .security/semgrep.yml --error --metrics=off --exclude .git --exclude .worktrees --exclude node_modules .; else docker run --rm -v "$ROOT:/src:ro" semgrep/semgrep:1.136.0 semgrep scan --config /src/.security/semgrep.yml --error --metrics=off --exclude .git --exclude .worktrees --exclude node_modules /src; fi; }
run_trivy(){ if have trivy; then trivy fs --scanners misconfig,secret --severity HIGH,CRITICAL --exit-code 1 --skip-dirs .git,.worktrees,node_modules .; else docker run --rm -v "$ROOT:/repo:ro" aquasec/trivy:0.66.0 fs --scanners misconfig,secret --severity HIGH,CRITICAL --exit-code 1 --skip-dirs /repo/.git,/repo/.worktrees,/repo/node_modules /repo; fi; }
case "${1:-all}" in all) run_gitleaks; run_semgrep; run_trivy;; gitleaks) run_gitleaks;; semgrep) run_semgrep;; trivy) run_trivy;; *) echo 'usage: run_security_scan.sh [all|gitleaks|semgrep|trivy]' >&2; exit 2;; esac
