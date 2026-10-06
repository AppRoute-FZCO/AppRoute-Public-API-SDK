/**
 * Product types in the catalog.
 */
export const ProductType = {
  VOUCHER: "voucher",
  DIRECT_TOPUP: "direct_topup",
  ESIM: "esim",
} as const;

export type ProductType = (typeof ProductType)[keyof typeof ProductType];
