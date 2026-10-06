import { BaseResource } from "./base-resource.js";
import type {
  ItemLookupRequestItem,
  ItemLookupResponse,
  Product,
  ProductItem,
  ProductListResponse,
  ProductStockResponse,
} from "../models/index.js";

/**
 * Batch lookup hard cap. Matches the backend's `max_length=100` on
 * `ItemLookupRequest.items`. Enforce client-side too so callers get a
 * fail-fast `Error` instead of a server-side 422 round-trip.
 */
export const MAX_LOOKUP_ITEMS = 100;

/**
 * Resource for the /services endpoints (product catalog).
 */
export class ServicesResource extends BaseResource {
  /**
   * List all products/services in the catalog.
   */
  async list(view?: string): Promise<ProductListResponse> {
    return this.transport.request<ProductListResponse>("GET", "/services", view === undefined ? undefined : { params: { view } });
  }

  /**
   * Get a single product/service by ID.
   */
  async get(productId: string, view?: string): Promise<Product> {
    return this.transport.request<Product>("GET", `/services/${productId}`, view === undefined ? undefined : { params: { view } });
  }

  /**
   * Get stock info for a product.
   */
  async stock(productId: string): Promise<ProductStockResponse> {
    return this.transport.request<ProductStockResponse>(
      "GET",
      `/services/${productId}/stock`,
    );
  }

  /**
   * Get a single denomination/item from a service by id.
   *
   * Calls `GET /services/{serviceId}/items/{itemId}` and returns the same
   * `ProductItem` shape that appears inside
   * `GET /services/{serviceId}.data.items[]`.
   */
  async getItem(serviceId: string, itemId: string, view?: string): Promise<ProductItem> {
    // Pre-existing convention: ids are server-issued UUIDs; not URI-encoded.
    // See README. Matches behaviour of `get(productId)` above.
    return this.transport.request<ProductItem>(
      "GET",
      `/services/${serviceId}/items/${itemId}`,
      view === undefined ? undefined : { params: { view } },
    );
  }

  /**
   * Batch lookup of up to 100 `(serviceId, itemId)` pairs in one round-trip.
   *
   * The response `items` are in the same order as the input — callers can
   * `zip()` request and response without re-keying.
   *
   * Throws an `Error` when the input is empty or longer than 100. These
   * checks happen before any HTTP call is made.
   */
  async lookupItems(
    items: ItemLookupRequestItem[],
    view?: string,
  ): Promise<ItemLookupResponse> {
    if (items.length === 0) {
      throw new Error("items must not be empty");
    }
    if (items.length > MAX_LOOKUP_ITEMS) {
      throw new Error(
        `items must contain at most ${MAX_LOOKUP_ITEMS} entries`,
      );
    }
    return this.transport.request<ItemLookupResponse>(
      "POST",
      "/services/items/lookup",
      view === undefined ? { body: { items } } : { params: { view }, body: { items } },
    );
  }
}
