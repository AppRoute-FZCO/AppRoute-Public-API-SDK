"""Release tags must agree with every committed package version before publish."""
from __future__ import annotations

import json
import subprocess
import sys
from pathlib import Path

import pytest

SCRIPT = Path(__file__).parents[2] / "scripts/check-release-version.py"


@pytest.fixture
def release_tree(tmp_path: Path) -> Path:
    files = {
        "python/approute/_version.py": '__version__ = "1.2.0"\n',
        "javascript/package.json": json.dumps({"version": "1.2.0"}),
        "javascript/package-lock.json": json.dumps({
            "version": "1.2.0", "packages": {"": {"version": "1.2.0"}},
        }),
        "php/composer.json": json.dumps({
            "name": "approute/public-api-sdk", "version": "1.2.0",
            "require": {"php": ">=8.1", "guzzlehttp/guzzle": "^7.5"},
        }),
        "composer.json": json.dumps({
            "name": "approute/public-api-sdk", "version": "1.2.0",
            "require": {"php": ">=8.1", "guzzlehttp/guzzle": "^7.5"},
        }),
        "golang/version.go": 'package approute\nconst Version = "1.2.0"\n',
    }
    for name, content in files.items():
        path = tmp_path / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content)
    return tmp_path


def check_version(root: Path, tag: str) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [sys.executable, str(SCRIPT), tag, "--root", str(root)],
        capture_output=True, text=True, check=False,
    )


def test_release_accepts_agreeing_committed_versions(release_tree: Path) -> None:
    result = check_version(release_tree, "v1.2.0")
    assert result.returncode == 0, result.stderr


@pytest.mark.parametrize("path", [
    "python/approute/_version.py", "javascript/package.json",
    "javascript/package-lock.json", "php/composer.json", "composer.json", "golang/version.go",
])
def test_release_rejects_stale_package_version(release_tree: Path, path: str) -> None:
    file = release_tree / path
    file.write_text(file.read_text().replace("1.2.0", "1.1.0"))
    result = check_version(release_tree, "v1.2.0")
    assert result.returncode == 1
    assert "1.1.0" in result.stderr
    assert "1.2.0" in result.stderr
    assert "1.1.0" in file.read_text()  # Verification must not rewrite stale versions.


def test_release_rejects_stale_npm_root_lock_entry(release_tree: Path) -> None:
    file = release_tree / "javascript/package-lock.json"
    metadata = json.loads(file.read_text())
    metadata["packages"][""]["version"] = "1.1.0"
    file.write_text(json.dumps(metadata))
    result = check_version(release_tree, "v1.2.0")
    assert result.returncode == 1
    assert "1.1.0" in result.stderr


def test_release_rejects_different_root_composer_requirements(release_tree: Path) -> None:
    file = release_tree / "composer.json"
    metadata = json.loads(file.read_text())
    metadata["require"]["guzzlehttp/guzzle"] = "^7.6"
    file.write_text(json.dumps(metadata))
    result = check_version(release_tree, "v1.2.0")
    assert result.returncode == 1
    assert "require" in result.stderr


@pytest.mark.parametrize("tag", ["1.2.0", "v", "v1.2.0;publish", "v2.0.0"])
def test_release_rejects_invalid_or_mismatched_tag(release_tree: Path, tag: str) -> None:
    result = check_version(release_tree, tag)
    assert result.returncode == 1
