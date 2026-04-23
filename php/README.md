# AppRoute Public API SDK for PHP

Official PHP SDK for [AppRoute Public API](https://api.approute.io).

## Installation

```bash
composer require approute/public-api-sdk
```

## Quick Start

```php
<?php

use AppRoute\Sdk\AppRouteClient;

$client = new AppRouteClient(apiKey: 'sk_live_...');

// List products
$products = $client->services->list();
foreach ($products as $product) {
    echo "{$product->name} ({$product->type})\n";
}

// Get a product
$product = $client->services->get('product-uuid');

// Create an order
$order = $client->orders->create([
    'ordersType' => 'shop',
    'reference' => 'my-unique-ref',
    'itemId' => 'item-uuid',
    'quantity' => 1,
]);
echo "Order: {$order->status}\n";

// Check balances
$accounts = $client->accounts->balances();
foreach ($accounts as $acc) {
    echo "{$acc->currency}: {$acc->available}\n";
}

// Fund your account
$methods = $client->funds->methods();
$invoice = $client->funds->createInvoice('USDT_TRC20', 100.0);
echo "Send {$invoice->amountExpected} to {$invoice->address}\n";

// Steam currency rates
$rates = $client->steamCurrency->rates(quotes: ['RUB', 'KZT']);
```

## Configuration

```php
$client = new AppRouteClient(
    apiKey: 'sk_live_...',
    baseUrl: 'https://api.approute.io/api/v1',  // default
    timeout: 30.0,       // seconds
    maxRetries: 3,       // retries on 429/5xx
);
```

## Error Handling

```php
use AppRoute\Sdk\AppRouteClient;
use AppRoute\Sdk\Exceptions\NotFoundException;
use AppRoute\Sdk\Exceptions\RateLimitedException;
use AppRoute\Sdk\Exceptions\ApiException;

$client = new AppRouteClient(apiKey: 'sk_live_...');

try {
    $product = $client->services->get('nonexistent-id');
} catch (NotFoundException $e) {
    echo "Not found: {$e->getMessage()}\n";
    echo "Trace ID: {$e->traceId}\n";
} catch (RateLimitedException $e) {
    echo "Rate limited\n";
} catch (ApiException $e) {
    echo "API error [{$e->errorCode}]: {$e->getMessage()}\n";
    foreach ($e->errors as $err) {
        echo "  {$err['field']}: {$err['message']}\n";
    }
}
```

### Exception Hierarchy

- `AppRouteException` — base for all SDK exceptions
  - `NetworkException` — connection/timeout errors
  - `ApiException` — API returned an error
    - `ValidationException` — 422
    - `UnauthorizedException` — 401
    - `ForbiddenException` — 403
    - `NotFoundException` — 404
    - `ConflictException` — 409
    - `RateLimitedException` — 429
    - `OutOfStockException` — 422
    - `InsufficientFundsException` — 422
    - `UpstreamException` — 502
    - `InternalException` — 500

## Requirements

- PHP 8.1+
- guzzlehttp/guzzle ^7.5
