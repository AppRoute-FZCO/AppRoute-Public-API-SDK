<?php

declare(strict_types=1);

namespace AppRoute\Sdk\Http;

use AppRoute\Sdk\Contracts\TransportInterface;
use AppRoute\Sdk\Exceptions\ApiException;
use AppRoute\Sdk\Exceptions\AppRouteException;
use AppRoute\Sdk\Exceptions\NetworkException;
use GuzzleHttp\Client;
use GuzzleHttp\Exception\ConnectException;
use GuzzleHttp\Exception\RequestException;

class HttpTransport implements TransportInterface
{
    private const SUCCESS_CODES = ['OK', 'ACCEPTED', 'IDEMPOTENCY_REPLAY'];
    private const STATUS_CODES = [
        0 => 'OK',
        1 => 'ACCEPTED',
        2 => 'IDEMPOTENCY_REPLAY',
        3 => 'VALIDATION_ERROR',
        4 => 'UNAUTHORIZED',
        5 => 'FORBIDDEN',
        6 => 'NOT_FOUND',
        7 => 'CONFLICT',
        8 => 'LIMIT_REACHED',
        9 => 'OUT_OF_STOCK',
        10 => 'INSUFFICIENT_FUNDS',
        11 => 'UPSTREAM_ERROR',
        12 => 'INTERNAL_ERROR',
        13 => 'API_KEY_LIMIT_EXCEEDED',
        14 => 'PROVIDER_UNAVAILABLE',
    ];
    private const RETRYABLE_STATUS_CODES = [429, 500, 502, 503, 504];
    private const BACKOFF_BASE = 1.0;

    private Client $client;

    public function __construct(
        private readonly string $baseUrl,
        private readonly string $apiKey,
        float $timeout = 30.0,
        private readonly int $maxRetries = 3,
    ) {
        $this->client = new Client([
            'base_uri' => rtrim($baseUrl, '/') . '/',
            'timeout' => $timeout,
            'headers' => [
                'X-API-Key' => $apiKey,
                'Content-Type' => 'application/json',
                'Accept' => 'application/json',
            ],
        ]);
    }

    /**
     * @param array<string, mixed>|null $params
     * @param array<string, mixed>|null $body
     * @return array<string, mixed>|null
     */
    public function request(string $method, string $path, ?array $params = null, ?array $body = null): ?array
    {
        $path = ltrim($path, '/');
        $options = [];

        if ($params !== null) {
            $cleaned = $this->cleanParams($params);
            if ($cleaned) {
                $options['query'] = $cleaned;
            }
        }

        if ($body !== null) {
            $options['json'] = $body;
        }

        $lastException = null;

        for ($attempt = 0; $attempt <= $this->maxRetries; $attempt++) {
            try {
                $response = $this->client->request($method, $path, $options);
            } catch (ConnectException $e) {
                throw new NetworkException("Connection error: {$e->getMessage()}", 0, $e);
            } catch (RequestException $e) {
                if ($e->hasResponse()) {
                    $response = $e->getResponse();
                } else {
                    throw new NetworkException("Request error: {$e->getMessage()}", 0, $e);
                }
            }

            $statusCode = $response->getStatusCode();

            if (in_array($statusCode, self::RETRYABLE_STATUS_CODES) && $attempt < $this->maxRetries) {
                $retryAfter = $response->getHeaderLine('Retry-After');
                $delay = $retryAfter !== '' ? (float) $retryAfter : self::BACKOFF_BASE * (2 ** $attempt);
                usleep((int) ($delay * 1_000_000));
                continue;
            }

            return $this->handleResponse($response->getBody()->getContents(), $statusCode);
        }

        throw new AppRouteException('Max retries exceeded');
    }

    /**
     * @return array<string, mixed>|null
     */
    private function handleResponse(string $rawBody, int $statusCode): ?array
    {
        $body = json_decode($rawBody, true);

        if (json_last_error() !== JSON_ERROR_NONE) {
            throw new AppRouteException("Invalid JSON response (HTTP {$statusCode})");
        }
        if (!is_array($body) || array_is_list($body)) {
            throw new AppRouteException("Invalid API response envelope (HTTP {$statusCode})");
        }

        $apiStatusCode = null;
        if (array_key_exists('statusCode', $body)) {
            if (!is_int($body['statusCode'])) {
                throw new AppRouteException("Invalid API response envelope (HTTP {$statusCode})");
            }
            $apiStatusCode = $body['statusCode'];
            $code = self::STATUS_CODES[$apiStatusCode] ?? (string) $apiStatusCode;
            $message = $body['statusMessage'] ?? '';
        } else {
            $code = $body['code'] ?? null;
            $message = $body['message'] ?? '';
            if (!is_string($code) || $code === '') {
                throw new AppRouteException("Invalid API response envelope (HTTP {$statusCode})");
            }
        }
        $traceId = $body['traceId'] ?? '';
        $apiErrorCode = $body['errorCode'] ?? null;
        $rawErrors = $body['errors'] ?? [];

        if ($statusCode >= 200 && $statusCode < 300 && in_array($code, self::SUCCESS_CODES, true) && !$rawErrors && !$apiErrorCode) {
            return $body['data'] ?? null;
        }

        $errors = [];
        foreach ($rawErrors as $err) {
            $errors[] = [
                'field' => $err['field'] ?? '',
                'code' => $err['code'] ?? '',
                'message' => $err['message'] ?? '',
            ];
        }

        throw ApiException::fromCode($code, $message, $traceId, $statusCode, $errors, $apiStatusCode, $apiErrorCode);
    }

    /**
     * @param array<string, mixed> $params
     * @return array<string, mixed>
     */
    private function cleanParams(array $params): array
    {
        $cleaned = [];
        foreach ($params as $key => $value) {
            if ($value === null) {
                continue;
            }
            if (is_bool($value)) {
                $cleaned[$key] = $value ? 'true' : 'false';
            } else {
                $cleaned[$key] = $value;
            }
        }
        return $cleaned;
    }
}
