"""The local/CI scanner gate must preserve tool failures and require reports."""

from __future__ import annotations

import os
import shutil
import subprocess
from pathlib import Path

import pytest

SCRIPT = Path(__file__).resolve().parents[2] / "scripts/security/run_security_scan.sh"


@pytest.fixture
def scanner(tmp_path: Path):
    script = tmp_path / "scripts/security/run_security_scan.sh"
    script.parent.mkdir(parents=True)
    shutil.copyfile(SCRIPT, script)
    binary = tmp_path / "bin"
    binary.mkdir()
    docker = binary / "docker"
    docker.write_text(
        "#!/usr/bin/env python3\n"
        "import os, pathlib, sys\n"
        "args = sys.argv[1:]\n"
        "pathlib.Path(os.environ['TRACE']).write_text(' '.join(args))\n"
        "if os.environ.get('TOOL_EXIT'): sys.exit(int(os.environ['TOOL_EXIT']))\n"
        "if '--output' in args and os.environ.get('WRITE_REPORT') == 'yes':\n"
        "    mount = next(a for a in args if a.endswith(':/reports'))\n"
        "    name = pathlib.Path(args[args.index('--output') + 1]).name\n"
        "    (pathlib.Path(mount[:-9]) / name).write_text('{}')\n"
    )
    docker.chmod(0o755)
    environment = {
        **os.environ,
        "PATH": f"{binary}:{os.environ['PATH']}",
        "TRACE": str(tmp_path / "trace"),
    }

    def run(mode: str, **extra: str):
        return subprocess.run(  # noqa: S603 - fixed local fixture script; no shell interpolation
            ["/bin/bash", str(script), mode],
            env={**environment, **extra},
            capture_output=True,
            text=True,
            check=False,
        )

    return run, tmp_path


@pytest.mark.parametrize("mode", ["gitleaks", "semgrep", "trivy-config", "trivy-vuln"])
def test_tool_failure_cannot_be_reported_as_pass(scanner, mode):
    run, _ = scanner
    assert run(mode, TOOL_EXIT="23").returncode == 23


def test_vulnerability_report_requires_actual_output(scanner):
    run, root = scanner
    assert run("trivy-vuln").returncode != 0
    assert run("trivy-vuln", WRITE_REPORT="yes").returncode == 0
    assert (root / ".security/results/trivy-vuln.json").is_file()


def test_secrets_are_redacted_and_source_is_readonly(scanner):
    run, root = scanner
    assert run("gitleaks").returncode == 0
    invocation = (root / "trace").read_text()
    assert "--redact=100" in invocation
    assert ":/repo:ro" in invocation
    assert "--network none" in invocation


def test_unknown_mode_never_invokes_scanner(scanner):
    run, root = scanner
    assert run("no-such-scan").returncode == 2
    assert not (root / "trace").exists()
