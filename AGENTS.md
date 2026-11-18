# Public API SDKs

This repository ships Python, JavaScript, PHP and Go clients for the same public API. When changing a public operation or response envelope, compare its OpenAPI definition with each affected client and the shared response fixtures.

Use the per-language commands in the CI workflow. For a published version change, follow the release procedure and keep the package metadata in sync.

- [OpenAPI contract](openapi.yaml)
- [Shared response examples](testdata/envelopes.json)
- [Language test commands](.github/workflows/ci.yml)
- [Release procedure](RELEASING.md)
