<?php

declare(strict_types=1);

namespace AppRoute\Sdk\Tests\Unit;

use AppRoute\Sdk\Exceptions\ApiException;
use AppRoute\Sdk\Exceptions\AppRouteException;
use AppRoute\Sdk\Http\HttpTransport;
use GuzzleHttp\Client;
use GuzzleHttp\Handler\MockHandler;
use GuzzleHttp\HandlerStack;
use GuzzleHttp\Psr7\Response;
use PHPUnit\Framework\Attributes\DataProvider;
use PHPUnit\Framework\TestCase;
use ReflectionClass;

final class EnvelopeContractTest extends TestCase
{
    public static function wireEnvelopes(): iterable
    {
        $fixtures = json_decode(file_get_contents(__DIR__ . '/../../../testdata/envelopes.json'), true, flags: JSON_THROW_ON_ERROR);
        foreach ($fixtures as $fixture) {
            yield $fixture['name'] => [$fixture];
        }
    }

    #[DataProvider('wireEnvelopes')]
    public function testPublicApiWireEnvelope(array $fixture): void
    {
        $handler = HandlerStack::create(new MockHandler([
            new Response($fixture['httpStatus'], [], json_encode($fixture['body'], JSON_THROW_ON_ERROR)),
        ]));
        $transport = new HttpTransport('https://api.test/api/v1', 'key', maxRetries: 0);
        (new ReflectionClass($transport))->getProperty('client')->setValue($transport, new Client(['handler' => $handler]));

        if ($fixture['invalid'] ?? false) {
            $this->expectException(AppRouteException::class);
            $this->expectExceptionMessage('envelope');
            $transport->request('GET', '/services');
        } elseif (array_key_exists('data', $fixture)) {
            $this->assertSame($fixture['data'], $transport->request('GET', '/services'));
        } else {
            $class = 'AppRoute\\Sdk\\Exceptions\\' . $fixture['errorType'] . 'Exception';
            try {
                $transport->request('GET', '/services');
                $this->fail('Expected API error');
            } catch (ApiException $error) {
                $this->assertInstanceOf($class, $error);
                $this->assertSame($fixture['code'], $error->errorCode);
                $this->assertStringContainsString($fixture['message'], $error->getMessage());
                $this->assertSame($fixture['body']['traceId'], $error->traceId);
                $this->assertSame($fixture['httpStatus'], $error->statusCode);
                $this->assertSame($fixture['apiStatusCode'] ?? null, $error->apiStatusCode);
                $this->assertSame($fixture['apiErrorCode'] ?? null, $error->apiErrorCode);
                $this->assertSame($fixture['errors'] ?? [], $error->errors);
            }
        }
    }
}
