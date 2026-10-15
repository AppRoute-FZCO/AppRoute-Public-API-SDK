# SDK release configuration

All four SDKs use the same committed version. Run the repository bump helper,
then `python scripts/check-release-version.py v1.2.0` (using the intended tag).
The release workflow verifies all package versions, both npm lock entries and
the two Composer manifests before any publication job starts. Re-running npm's
version command with the already committed version is supported.

Merging a PR does not publish packages. The release workflow starts from a root
`v<version>` tag. Package publication, tag creation and registry configuration
are separate release operations and require successful CI plus configured
destinations and credentials.

## Package sources

Go's module path is
`github.com/AppRoute-FZCO/AppRoute-Public-API-SDK/golang`, matching the verified
public repository and its module subdirectory. Go needs the additional
`golang/v<version>` tag created by the release workflow. A root tag alone does
not make the submodule version fetchable. Update consumers that still import
the previously advertised `github.com/approute/public-api-sdk-go`; that source
was unavailable during the release audit. No redirect or compatible alias is
assumed.

Composer VCS discovery reads `composer.json` at the source repository root.
The root manifest declares `approute/public-api-sdk` and maps its namespace to
`php/src/`. The `php/composer.json` manifest is for developing the standalone
subtree. The bump helper updates both; release validation requires matching
package names, versions and runtime requirements. CI installs the root package
and instantiates the SDK with its generated autoloader, then runs the separate
subtree PHPUnit suite.

The registry operator must confirm that Satis indexes the VCS source
`https://github.com/AppRoute-FZCO/AppRoute-Public-API-SDK.git`, discovers the root
manifest at the intended tag and serves the resulting package metadata and
archives. The repository does not contain the external Satis configuration.
An SSH rebuild finishing is insufficient evidence that the expected package
version is indexed and installable.

## External publication prerequisites

The workflow expects these existing GitHub secret names:

| Destination | Configuration |
| --- | --- |
| Private PyPI | `REGISTRY_HOST`, `REGISTRY_USER`, `REGISTRY_PASS` |
| Private npm / Verdaccio | `REGISTRY_HOST`, `NPM_TOKEN` |
| Satis rebuild host | `DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_SSH_KEY` |

The repository secret/variable lists were empty in the 2026-10-15 audit.
Inherited organization configuration and actual registry credentials were not
verified. Configure or verify those values and the real registry destinations
before release; do not infer successful publication from green unit CI.

After publication, verify installation of each exact version from a fresh
consumer, including Go's submodule version and Composer's root autoload map.
The catalog rollout remains Shop first, then SDK publication and affected
consumer migration to explicit `view=full`, then Public API's compact default.
