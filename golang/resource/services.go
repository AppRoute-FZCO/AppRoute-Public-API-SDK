package resource

import (
	"context"

	"github.com/approute/public-api-sdk-go/model"
)

// ServicesResource provides access to the product/service catalog endpoints.
type ServicesResource struct {
	t Transport
}

// NewServices creates a new ServicesResource using the given transport.
func NewServices(t Transport) *ServicesResource {
	return &ServicesResource{t: t}
}

// List returns all products in the catalog.
func (r *ServicesResource) List(ctx context.Context) (*model.ProductListResponse, error) {
	raw, err := r.t.Request(ctx, "GET", "/services", nil, nil)
	if err != nil {
		return nil, err
	}
	return decode[model.ProductListResponse](raw, "ProductListResponse")
}

// Get retrieves a single product by its ID.
func (r *ServicesResource) Get(ctx context.Context, productID string) (*model.Product, error) {
	raw, err := r.t.Request(ctx, "GET", "/services/"+productID, nil, nil)
	if err != nil {
		return nil, err
	}
	return decode[model.Product](raw, "Product")
}

// Stock returns the current stock levels for a product.
func (r *ServicesResource) Stock(ctx context.Context, productID string) (*model.ProductStockResponse, error) {
	raw, err := r.t.Request(ctx, "GET", "/services/"+productID+"/stock", nil, nil)
	if err != nil {
		return nil, err
	}
	return decode[model.ProductStockResponse](raw, "ProductStockResponse")
}
