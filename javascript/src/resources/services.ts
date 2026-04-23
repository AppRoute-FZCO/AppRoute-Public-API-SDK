import { BaseResource } from "./base-resource.js";
import type {
  Product,
  ProductListResponse,
  ProductStockResponse,
} from "../models/index.js";

/**
 * Resource for the /services endpoints (product catalog).
 */
export class ServicesResource extends BaseResource {
  /**
   * List all products/services in the catalog.
   */
  async list(): Promise<ProductListResponse> {
    return this.transport.request<ProductListResponse>("GET", "/services");
  }

  /**
   * Get a single product/service by ID.
   */
  async get(productId: string): Promise<Product> {
    return this.transport.request<Product>("GET", `/services/${productId}`);
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
}
