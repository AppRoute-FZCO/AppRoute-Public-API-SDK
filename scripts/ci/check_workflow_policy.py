#!/usr/bin/env python3
"""Validate GitHub Actions policy without executing untrusted workflow code."""

from __future__ import annotations

import re
import tomllib
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
WORKFLOWS = ROOT / ".github" / "workflows"
EXCLUSIONS = ROOT / ".security" / "workflow-policy-exclusions.toml"
SHA_REF = re.compile(r"^[0-9a-f]{40}(?:\s+#\s+.+)?$")
USES = re.compile(r"^\s*-?\s*uses:\s*([^\s]+)", re.MULTILINE)

errors: list[str] = []
exclusions: dict[str, str] = {}
if EXCLUSIONS.exists():
    for item in tomllib.loads(EXCLUSIONS.read_text(encoding="utf-8")).get("exclude", []):
        path, reason = item.get("path"), item.get("reason", "").strip()
        if not path or not reason:
            errors.append(f"{EXCLUSIONS}: every exclusion needs path and reason")
        exclusions[path] = reason
for workflow in sorted((*WORKFLOWS.glob("*.yml"), *WORKFLOWS.glob("*.yaml"))):
    rel = workflow.relative_to(ROOT).as_posix()
    if rel in exclusions:
        continue
    text = workflow.read_text(encoding="utf-8")
    if "runs-on: ubuntu-latest" in text:
        errors.append(f"{workflow}: moving ubuntu-latest runner is forbidden")
    # Self-hosted runners are approved for this repository. Workflows must
    # retain the Linux label so Windows/macOS runners cannot pick them up.
    if "runs-on: self-hosted" in text:
        errors.append(f"{workflow}: scalar self-hosted runner labels are forbidden")
    if "runs-on: [self-hosted" in text and "linux" not in text:
        errors.append(f"{workflow}: self-hosted jobs must include the linux label")
    for action in USES.findall(text):
        if action.startswith(("./", "docker://")):
            continue
        if "@" not in action:
            errors.append(f"{workflow}: action without revision: {action}")
            continue
        _name, revision = action.rsplit("@", 1)
        if not SHA_REF.fullmatch(revision):
            errors.append(f"{workflow}: action is not pinned to a commit SHA: {action}")

if errors:
    raise SystemExit("\n".join(errors))
print(f"validated {len(tuple(WORKFLOWS.glob('*.yml')))} workflows")
