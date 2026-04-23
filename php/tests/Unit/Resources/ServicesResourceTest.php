<?php

declare(strict_types=1);

namespace AppRoute\Sdk\Tests\Unit\Resources;

use AppRoute\Sdk\Models\Product;
use AppRoute\Sdk\Models\ProductField;
use AppRoute\Sdk\Models\ProductFieldValidation;
use AppRoute\Sdk\Models\ProductItem;
use AppRoute\Sdk\Models\ProductStockItem;
use AppRoute\Sdk\Models\ProductStockResponse;
use AppRoute\Sdk\Resources\ServicesResource;
use AppRoute\Sdk\Tests\Factory\ProductFactory;
use AppRoute\Sdk\Tests\TestCase;

final class ServicesResourceTest extends TestCase
{
    private ServicesResource $resource;

    protected function setUp(): void
    {
        parent::setUp();
        $this->resource = new ServicesResource($this->transport);
    }

    // --- list() ---

    public function testListReturnsProductArray(): void
    {
        $this->transport->setResponse(ProductFactory::makeProductListData());

        $result = $this->resource->list();

        $this->assertCount(1, $result);
        $this->assertInstanceOf(Product::class, $result[0]);
        $this->assertSame('prod-001', $result[0]->id);
        $this->assertSame('Steam Wallet 50 USD', $result[0]->name);
        $this->assertSame('voucher', $result[0]->type);
    }

    public function testListCallsCorrectEndpoint(): void
    {
        $this->transport->setResponse(ProductFactory::makeProductListData());

        $this->resource->list();

        $call = $this->transport->lastCall();
        $this->assertSame('GET', $call['method']);
        $this->assertSame('/services', $call['path']);
        $this->assertNull($call['params']);
        $this->assertNull($call['body']);
    }

    public function testListProductContainsTypedItems(): void
    {
        $this->transport->setResponse(ProductFactory::makeProductListData());

        $products = $this->resource->list();

        $this->assertCount(1, $products[0]->items);
        $this->assertInstanceOf(ProductItem::class, $products[0]->items[0]);
        $this->assertSame('item-001', $products[0]->items[0]->id);
        $this->assertSame(48.5, $products[0]->items[0]->price);
    }

    public function testListProductContainsTypedFields(): void
    {
        $this->transport->setResponse(ProductFactory::makeProductListData());

        $products = $this->resource->list();

        $this->assertNotNull($products[0]->fields);
        $this->assertCount(1, $products[0]->fields);
        $this->assertInstanceOf(ProductField::class, $products[0]->fields[0]);
        $this->assertSame('email', $products[0]->fields[0]->key);
        $this->assertTrue($products[0]->fields[0]->required);
    }

    public function testListFieldHasTypedValidation(): void
    {
        $this->transport->setResponse(ProductFactory::makeProductListData());

        $products = $this->resource->list();
        $field = $products[0]->fields[0];

        $this->assertInstanceOf(ProductFieldValidation::class, $field->validation);
        $this->assertSame('^[^@]+@[^@]+$', $field->validation->pattern);
        $this->assertSame('Invalid email', $field->validation->message);
    }

    // --- get() ---

    public function testGetReturnsProduct(): void
    {
        $this->transport->setResponse(ProductFactory::makeProductGetData());

        $product = $this->resource->get('prod-001');

        $this->assertInstanceOf(Product::class, $product);
        $this->assertSame('prod-001', $product->id);
        $this->assertNull($product->fields);
    }

    public function testGetCallsCorrectEndpoint(): void
    {
        $this->transport->setResponse(ProductFactory::makeProductGetData());

        $this->resource->get('prod-001');

        $call = $this->transport->lastCall();
        $this->assertSame('GET', $call['method']);
        $this->assertSame('/services/prod-001', $call['path']);
    }

    // --- stock() ---

    public function testStockReturnsTypedResponse(): void
    {
        $this->transport->setResponse(ProductFactory::makeProductStockData());

        $result = $this->resource->stock('prod-001');

        $this->assertInstanceOf(ProductStockResponse::class, $result);
        $this->assertSame('prod-001', $result->productId);
        $this->assertCount(2, $result->items);
    }

    public function testStockItemsAreTyped(): void
    {
        $this->transport->setResponse(ProductFactory::makeProductStockData());

        $result = $this->resource->stock('prod-001');

        $this->assertInstanceOf(ProductStockItem::class, $result->items[0]);
        $this->assertSame('item-001', $result->items[0]->itemId);
        $this->assertSame(150, $result->items[0]->stock);
        $this->assertNull($result->items[1]->stock);
    }

    public function testStockCallsCorrectEndpoint(): void
    {
        $this->transport->setResponse(ProductFactory::makeProductStockData());

        $this->resource->stock('prod-001');

        $call = $this->transport->lastCall();
        $this->assertSame('GET', $call['method']);
        $this->assertSame('/services/prod-001/stock', $call['path']);
    }
}
