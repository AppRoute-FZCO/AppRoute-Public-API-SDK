# AppRoute Public API SDK

Official SDK libraries for [AppRoute Public API](https://api.approute.io) — available in Python, JavaScript/TypeScript, PHP, and Go.

## Available SDKs

| Language | Package | Min Version |
|----------|---------|-------------|
| [Python](./python/) | `approute-public-api-sdk` | Python 3.10+ |
| [JavaScript/TypeScript](./javascript/) | `@approute/public-api-sdk` | Node.js 18+ |
| [PHP](./php/) | `approute/public-api-sdk` | PHP 8.1+ |
| [Go](./golang/) | `github.com/approute/public-api-sdk-go` | Go 1.21+ |

## Quick Start

All SDKs follow the same pattern:

1. Get your API key from the [AppRoute Dashboard](https://lk.approute.io)
2. Install the SDK for your language
3. Initialize the client with your API key
4. Call API methods

## API Coverage

All SDKs cover the **Data Plane** endpoints:

- **Services** — product catalog (list, get by ID)
- **Orders** — create purchases, DTU checks, list orders
- **Accounts** — balances, transaction history
- **Funds** — funding methods, invoices (create, list, check), TON deposits, Bybit UID management
- **Steam Currency** — exchange rates

## Authentication

All endpoints use API key authentication via `X-API-Key` header. The SDK handles this automatically.

## Response Handling

All API responses use a unified envelope format. The SDKs automatically:
- Unwrap successful responses, returning just the `data` field
- Throw/return typed errors for failed responses with `code`, `message`, `traceId`, and field-level `errors`

## License

MIT
