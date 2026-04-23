<?php

declare(strict_types=1);

namespace AppRoute\Sdk\Resources;

use AppRoute\Sdk\Contracts\ResourceInterface;
use AppRoute\Sdk\Contracts\TransportInterface;
use AppRoute\Sdk\Models\Product;
use AppRoute\Sdk\Models\ProductStockResponse;

class ServicesResource implements ResourceInterface
{
    public function __construct(private readonly TransportInterface $transport) {}

    /** @return Product[] */
    public function list(): array
    {
        $data = $this->transport->request('GET', '/services');
        return array_map(fn(array $p) => Product::fromArray($p), $data['items'] ?? []);
    }

    public function get(string $productId): Product
    {
        $data = $this->transport->request('GET', "/services/{$productId}");
        return Product::fromArray($data);
    }

    public function stock(string $productId): ProductStockResponse
    {
        $data = $this->transport->request('GET', "/services/{$productId}/stock");
        return ProductStockResponse::fromArray($data);
    }
}
