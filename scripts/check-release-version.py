#!/usr/bin/env python3
"""Verify committed SDK versions before any release job publishes artifacts."""

from __future__ import annotations

import argparse
import ast
import json
import re
import sys
from pathlib import Path


def check(root: Path, tag: str) -> None:
    match = re.fullmatch(r"v([0-9]+\.[0-9]+\.[0-9]+(?:-[a-zA-Z0-9.]+)?)", tag)
    if match is None:
        raise ValueError("Release tag must be v<semver>, e.g. v1.2.0")
    expected = match.group(1)

    python_tree = ast.parse((root / "python/approute/_version.py").read_text())
    python_version = next(
        ast.literal_eval(node.value)
        for node in python_tree.body
        if isinstance(node, ast.Assign)
        and any(
            isinstance(target, ast.Name) and target.id == "__version__"
            for target in node.targets
        )
    )
    javascript = json.loads((root / "javascript/package.json").read_text())
    lock = json.loads((root / "javascript/package-lock.json").read_text())
    php = json.loads((root / "php/composer.json").read_text())
    composer = json.loads((root / "composer.json").read_text())
    go_match = re.search(
        r'^const Version = "([^"]+)"$',
        (root / "golang/version.go").read_text(),
        re.MULTILINE,
    )
    if go_match is None:
        raise ValueError("Cannot read exported Go Version")
    versions = {
        "Python": python_version,
        "JavaScript": javascript["version"],
        "npm lock": lock["version"],
        "npm root lock entry": lock["packages"][""]["version"],
        "PHP subtree": php["version"],
        "PHP distribution": composer["version"],
        "Go": go_match.group(1),
    }
    failures = [
        f"{name}: committed {version}, tag requires {expected}"
        for name, version in versions.items()
        if version != expected
    ]
    for key in ("name", "require"):
        if php[key] != composer[key]:
            failures.append(f"Root and subtree composer.json {key} must agree")
    if failures:
        raise ValueError("\n".join(failures))
    print(f"All four SDKs, npm lock entries and Composer manifests agree with {tag}")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("tag")
    parser.add_argument(
        "--root", type=Path, default=Path(__file__).resolve().parents[1]
    )
    args = parser.parse_args()
    try:
        check(args.root, args.tag)
    except (OSError, ValueError, KeyError, SyntaxError, StopIteration) as exc:
        print(f"Release verification failed: {exc}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
