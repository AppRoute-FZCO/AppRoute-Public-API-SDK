import { describe, it, expect } from "vitest";
import { ServicesResource } from "../../src/resources/services.js";
import { MockTransport } from "../support/index.js";
import { NotFoundError } from "../../src/errors/index.js";
import {
  makeProductListResponse,
  makeProductGetResponse,
  makeProductStockResponse,
} from "../factories/index.js";
import type {
  ProductListResponse,
  Product,
  ProductStockResponse,
} from "../../src/models/index.js";

describe("ServicesResource", () => {
  describe("list()", () => {
    it("should return product list on success", async () => {
      const fixture = makeProductListResponse();
      const mock = new MockTransport(fixture);
      const resource = new ServicesResource(mock.transport);

      const result: ProductListResponse = await resource.list();

      expect(result.items).toHaveLength(1);
      expect(result.items[0].id).toBe("prod-001");
      expect(result.items[0].name).toBe("Steam Wallet 50 USD");
      expect(result.items[0].type).toBe("voucher");
      expect(result.items[0].items).toHaveLength(1);
      expect(result.items[0].items[0].price).toBe(48.5);
      expect(result.hasNext).toBe(false);
    });

    it("should call GET /services", async () => {
      const mock = new MockTransport(makeProductListResponse());
      const resource = new ServicesResource(mock.transport);

      await resource.list();

      expect(mock.callCount).toBe(1);
      expect(mock.lastCall.method).toBe("GET");
      expect(mock.lastCall.path).toBe("/services");
      expect(mock.lastCall.options).toBeUndefined();
    });

    it("should include fields in product list", async () => {
      const fixture = makeProductListResponse();
      const mock = new MockTransport(fixture);
      const resource = new ServicesResource(mock.transport);

      const result = await resource.list();
      const fields = result.items[0].fields;

      expect(fields).toHaveLength(1);
      expect(fields![0].key).toBe("email");
      expect(fields![0].type).toBe("text");
      expect(fields![0].required).toBe(true);
      expect(fields![0].validation?.pattern).toBe("^[^@]+@[^@]+$");
    });
  });

  describe("get()", () => {
    it("should return a single product on success", async () => {
      const fixture = makeProductGetResponse();
      const mock = new MockTransport(fixture);
      const resource = new ServicesResource(mock.transport);

      const result: Product = await resource.get("prod-001");

      expect(result.id).toBe("prod-001");
      expect(result.name).toBe("Steam Wallet 50 USD");
      expect(result.countryCode).toBe("US");
      expect(result.categoryName).toBe("Gaming");
    });

    it("should call GET /services/:id", async () => {
      const mock = new MockTransport(makeProductGetResponse());
      const resource = new ServicesResource(mock.transport);

      await resource.get("prod-001");

      expect(mock.callCount).toBe(1);
      expect(mock.lastCall.method).toBe("GET");
      expect(mock.lastCall.path).toBe("/services/prod-001");
    });

    it("should propagate errors from transport", async () => {
      const error = new NotFoundError("NOT_FOUND", "Product not found", "trace-1", 404);
      const mock = new MockTransport().setError(error);
      const resource = new ServicesResource(mock.transport);

      await expect(resource.get("unknown")).rejects.toThrow(NotFoundError);
    });
  });

  describe("stock()", () => {
    it("should return stock information on success", async () => {
      const fixture = makeProductStockResponse();
      const mock = new MockTransport(fixture);
      const resource = new ServicesResource(mock.transport);

      const result: ProductStockResponse = await resource.stock("prod-001");

      expect(result.productId).toBe("prod-001");
      expect(result.items).toHaveLength(2);
      expect(result.items[0].itemId).toBe("item-001");
      expect(result.items[0].stock).toBe(150);
      expect(result.items[1].stock).toBeNull();
    });

    it("should call GET /services/:id/stock", async () => {
      const mock = new MockTransport(makeProductStockResponse());
      const resource = new ServicesResource(mock.transport);

      await resource.stock("prod-001");

      expect(mock.callCount).toBe(1);
      expect(mock.lastCall.method).toBe("GET");
      expect(mock.lastCall.path).toBe("/services/prod-001/stock");
    });
  });
});
