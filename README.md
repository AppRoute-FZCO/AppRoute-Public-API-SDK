# AppRoute Public API SDK

Official SDK libraries for [AppRoute Public API](https://approute.io) — available in Python, JavaScript/TypeScript, PHP, and Go.

## Available SDKs

| Language | Package | Min Version |
|----------|---------|-------------|
| [Python](./python/) | `approute-public-api-sdk` | Python 3.10+ |
| [JavaScript/TypeScript](./javascript/) | `@approute/public-api-sdk` | Node.js 18+ |
| [PHP](./php/) | `approute/public-api-sdk` | PHP 8.1+ |
| [Go](./golang/) | `github.com/AppRoute-FZCO/AppRoute-Public-API-SDK/golang` | Go 1.24+ |

## Quick Start

All SDKs follow the same pattern:

1. Get your API key from the AppRoute Dashboard:
   - International clients: [approute.io/dashboard](https://approute.io/dashboard)
   - Russian clients: [approute.ru/dashboard](https://approute.ru/dashboard)
2. Install the SDK for your language
3. Initialize the client with your API key
4. Call API methods

## Regional Domains

Use the `.io` domain for international clients and the `.ru` domain for Russian clients:

| Region | Dashboard | API Base URL |
|--------|-----------|--------------|
| International | `https://approute.io/dashboard` | `https://approute.io/api/v1` |
| Russia | `https://approute.ru/dashboard` | `https://approute.ru/api/v1` |

The SDKs default to the international API endpoint. Russian clients should pass the `.ru` API base URL in client configuration.

## API Coverage

All SDKs cover the **Data Plane** endpoints:

- **Services** — product catalog (list, get by ID, get single item, batch item lookup)
- **Orders** — create purchases, DTU checks, list orders
- **Accounts** — balances, transaction history
- **Funds** — funding methods, invoices (create, list, check), TON deposits, Bybit UID management
- **Steam Currency** — exchange rates

### Catalog response views

Catalog methods accept `view="compact"` or `view="full"`; omitted `view` defaults
to compact. Use full during migration when clients need images and the other
allowlisted metadata:

```text
GET /services
→
GET /services?view=full
```

The same parameter is supported by the service, single-item, and batch-lookup
methods in every SDK. Unknown values are rejected by the API with HTTP 422.

Full product responses include available `imageUrl`, `section`, `categoryName`,
`subcategoryName`, `countryCode`, `activationCountries`, and `activationGuide`.
These keys are omitted when unavailable. eSIM images also require full. Single
items currently have no extra presentation metadata; their view parameter is
forwarded consistently. No item image or platform source is introduced here.

**Deployment migration:** the API's compact default changes the response for
existing callers, including older SDK versions, when the API is deployed.
Deploy Shop first. Release the compatible SDKs and migrate affected callers to explicit `view=full` before
deploying Public API's compact default. The previous Public API ignores this
query parameter, so callers can stage it in advance. SDK 1.2.0 adds optional
arguments and fields without changing existing call syntax; the Go release
uses `github.com/AppRoute-FZCO/AppRoute-Public-API-SDK/golang` and the
`golang/v1.2.0` tag. Go consumers must update imports from the previously
documented, unavailable `github.com/approute/public-api-sdk-go` source.

All four transports support the current numeric `statusCode`/`statusMessage`
envelope and the legacy `code`/`message` envelope. See [release configuration](./RELEASING.md)
for registry and package source requirements before publishing.

## Authentication

All endpoints use API key authentication via `X-API-Key` header. The SDK handles this automatically.

## Response Handling

All API responses use a unified envelope format. The SDKs automatically:
- Unwrap successful responses, returning just the `data` field
- Throw/return typed errors for failed responses with `code`, `message`, `traceId`, and field-level `errors`

Current `statusCode` values 0 (`OK`), 1 (`ACCEPTED`) and 2
(`IDEMPOTENCY_REPLAY`) succeed only with HTTP 2xx and no field or machine error.
A present numeric code takes precedence over legacy fields. Unknown numeric
codes remain API errors; malformed codes are rejected instead of coerced to 0.

Existing SDK error code and HTTP status properties keep their meanings. Errors
also retain the numeric API code as `apiStatusCode` and the wire `errorCode` as
`apiErrorCode` (`api_status_code`/`api_error_code` in Python,
`APIStatusCode`/`APIErrorCode` in Go). In PHP, `errorCode` remains the existing
symbolic result code. Code 13 maps to the conflict error type; code 14 maps to
the upstream error type. These codes are distinct from a rate limit (code 8).

## Changelog

### 1.2.0 (unreleased)

Adds compact/full catalog views in all four SDKs, optional product presentation
and quantity-mode metadata, the eSIM product type, and nullable
`FundingInvoice.explorerUrl`. Fixes current API response envelopes in all four
transports, aligns Go module discovery and adds a root Composer distribution
manifest. See the catalog and Go import migration notes above.

### 1.1.0

Added `ServicesResource.getItem` (single-item GET) and `lookupItems` (batch
POST, up to 100 `(serviceId, itemId)` pairs). No breaking changes.

### Breaking change (2026-05): `POST /orders` no longer returns 429 for API-key spend caps

**What changed.** `POST /api/v1/orders` previously returned HTTP 429
`LIMIT_REACHED` in two structurally different cases:

1. A genuine per-window rate-limit (Public-API or upstream shop-service).
2. The API key's static `transactionLimit` could not accommodate the order.

These are now disambiguated:

| Case | HTTP | `statusCode` | `errorCode` | Retryable? |
|------|------|--------------|-------------|------------|
| Per-window rate-limit | `429` | `LIMIT_REACHED` (8) | — | Yes — back-off retry |
| API-key spend cap | `409` | `API_KEY_LIMIT_EXCEEDED` (13) | `api_key_transaction_limit_exceeded` | **No** — raise the cap or reset `transactionUsed` |

**Migration.** If your application had special handling on HTTP 429
`LIMIT_REACHED` for the "API key transaction limit exhausted" case
(typically: stop retrying, surface to operator, page support), move that
branch onto HTTP 409 with
`errorCode === "api_key_transaction_limit_exceeded"`. Keep your existing
back-off path on HTTP 429 `LIMIT_REACHED` — it is now strictly the
per-window rate-limit signal.

The new error envelope's `statusMessage` for the spend-cap path includes
the attempted amount, the remaining headroom and the configured limit so
operators can act on it directly:

```
Order amount 0.5 USD exceeds api key remaining limit 0.2 of 1.0
```

The control-plane representation of an API key now also exposes
`transactionRemaining` (`transactionLimit - transactionUsed`, or `null`
when no limit is set). Use it instead of `status` to decide whether a
key can serve an order of expected size — a key with `status="active"`
can still refuse new orders with HTTP 409 if its remaining headroom is
below the quoted amount.

## License

MIT
