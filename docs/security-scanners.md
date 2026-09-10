# Static security checks

Run `bash scripts/security/run_security_scan.sh SCANNER`, where SCANNER is
gitleaks, semgrep, trivy-config or trivy-vuln. Local and CI use the same pinned
Linux/amd64 images. Secret output is redacted; source is mounted read-only.
Gitleaks and Semgrep run without network access. Semgrep uses local rules and
strict parsing without relying on a worktree gitdir inside the container.
The source snapshot is scanned; this is not a Git-history audit.

Gitleaks, Semgrep and high/critical configuration findings block the gate.
Dependency CVEs are report-only while scanner failures or absent reports fail.
The dependency report is retained for 14 days; inspect and remediate its
findings separately. Reports do not establish a comprehensive security review.

The workflow uses a hosted Ubuntu runner and read-only repository permissions.
It does not deploy, send notifications, access production, invoke licensed
CodeQL/dependency-review capabilities or run active network scans. The old ZAP
helper's arbitrary sandbox hostname/host-network behavior is excluded until an
isolated target and network policy are provisioned. Existing deploy runners and
service topology are unchanged. Original scanner branch evidence is retained.
