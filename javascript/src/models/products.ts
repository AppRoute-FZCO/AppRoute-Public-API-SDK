/**
 * A selectable option for a product field (e.g. dropdown values).
 */
export interface ProductFieldOption {
  label: string;
  value: string;
  price?: number;
}

/**
 * Validation rules for a product field.
 */
export interface ProductFieldValidation {
  min?: number;
  max?: number;
  pattern?: string;
  message?: string;
}

/**
 * A field definition on a product (input the buyer must fill).
 */
export interface ProductField {
  key: string;
  type: string;
  required: boolean;
  label?: string;
  options?: ProductFieldOption[];
  validation?: ProductFieldValidation;
}

/**
 * A purchasable item (denomination) within a product.
 */
export interface ProductItem {
  id: string;
  name?: string;
  nominal?: number;
  price: number;
  currency: string;
  available?: boolean;
  stock?: number;
  isLongOrder?: boolean;
  minQtyToLongOrder?: number;
}

/**
 * A product in the catalog.
 */
export interface Product {
  id: string;
  name?: string;
  type: string;
  imageUrl?: string;
  countryCode?: string;
  categoryName?: string;
  subcategoryName?: string;
  items: ProductItem[];
  fields?: ProductField[];
}

/**
 * Paginated response for product listings.
 */
export interface ProductListResponse {
  items: Product[];
  hasNext: boolean;
}

/**
 * Stock information for a single item within a product.
 */
export interface ProductStockItem {
  itemId: string;
  stock?: number;
}

/**
 * Stock information for all items within a product.
 */
export interface ProductStockResponse {
  productId: string;
  items: ProductStockItem[];
}
