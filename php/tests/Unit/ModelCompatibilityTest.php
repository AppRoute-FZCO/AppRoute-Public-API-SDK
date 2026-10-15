<?php

declare(strict_types=1);

namespace AppRoute\Sdk\Tests\Unit;

use AppRoute\Sdk\Models\FundingInvoice;
use AppRoute\Sdk\Models\Product;
use AppRoute\Sdk\Models\ProductItem;
use PHPUnit\Framework\TestCase;

final class ModelCompatibilityTest extends TestCase
{
    public function testFullCatalogMetadataFromWire(): void
    {
        $product = Product::fromArray([
            'id' => 'svc-1', 'name' => null, 'type' => 'esim', 'fields' => null,
            'imageUrl' => 'https://flagcdn.com/us.svg', 'section' => 'esim',
            'activationCountries' => 'US', 'activationGuide' => 'Scan QR code',
            'items' => [['id' => 'item-1', 'price' => 9.5, 'currency' => 'USD', 'inStock' => 10]],
        ]);

        $this->assertSame('https://flagcdn.com/us.svg', $product->imageUrl);
        $this->assertSame('US', $product->activationCountries);
        $this->assertSame('Scan QR code', $product->activationGuide);
        $this->assertSame('esim', $product->section);
        $this->assertSame(10, $product->items[0]->inStock);
        $this->assertNull($product->items[0]->nominal);
    }

    public function testProductItemPreservesExistingPositionalArguments(): void
    {
        $item = new ProductItem('item-1', 'Voucher', 10.0, 9.5, 'USD', true, 20, true, 5);

        $this->assertTrue($item->isLongOrder);
        $this->assertSame(5, $item->minQtyToLongOrder);
        $this->assertNull($item->inStock);
    }

    public function testFundingInvoicePreservesExistingPositionalArguments(): void
    {
        $invoice = new FundingInvoice(
            'inv-1', 'USDT_TRC20', 10.0, 0.0, 10.0, 'address', 'success',
            '2026-10-14T00:00:00Z', '2026-10-14T01:00:00Z', 'incoming', null,
            'hash', 20, 21,
        );

        $this->assertSame(20, $invoice->confirmationsRequired);
        $this->assertSame(21, $invoice->confirmations);
        $this->assertNull($invoice->explorerUrl);
    }
}
