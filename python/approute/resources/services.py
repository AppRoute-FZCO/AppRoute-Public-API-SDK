from approute.models import Product, ProductListResponse, ProductStockResponse
from approute.resources._base import BaseResource


class ServicesResource(BaseResource):

    def list(self) -> ProductListResponse:
        """List all products/services in the catalog."""
        data = self._t.request("GET", "/services")
        return ProductListResponse.model_validate(data)

    def get(self, product_id: str) -> Product:
        """Get a single product/service by ID."""
        data = self._t.request("GET", f"/services/{product_id}")
        return Product.model_validate(data)

    def stock(self, product_id: str) -> ProductStockResponse:
        """Get stock info for a product."""
        data = self._t.request("GET", f"/services/{product_id}/stock")
        return ProductStockResponse.model_validate(data)
