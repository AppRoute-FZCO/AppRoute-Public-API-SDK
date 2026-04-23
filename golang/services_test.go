package approute_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	approute "github.com/approute/public-api-sdk-go"
	"github.com/approute/public-api-sdk-go/apierror"
	"github.com/approute/public-api-sdk-go/enum"
	"github.com/approute/public-api-sdk-go/model"
)

func TestServicesResource_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "Method", r.Method, "GET")
		assertEqual(t, "Path", r.URL.Path, "/services")
		w.Header().Set("Content-Type", "application/json")
		w.Write(successEnvelope(t, newProductListResponse(
			newProduct(),
			newProduct(func(p *model.Product) {
				p.ID = "prod-2"
				p.Name = strPtr("Mobile Top-Up")
				p.Type = enum.ProductDirectTopup
				p.ImageURL = nil
				p.CountryCode = strPtr("TR")
				p.CategoryName = strPtr("Mobile")
				p.SubcategoryName = nil
				p.Items = []model.ProductItem{
					newProductItem(func(pi *model.ProductItem) {
						pi.ID = "item-2"
						pi.Name = nil
						pi.Nominal = 50.00
						pi.Price = 5.00
						pi.Stock = nil
					}),
				}
				p.Fields = []model.ProductField{newProductField()}
			}),
		)))
	}))
	defer srv.Close()

	client := approute.NewClient("test-key", approute.WithBaseURL(srv.URL), approute.WithMaxRetries(0))
	result, err := client.Services.List(t.Context())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 2 {
		t.Fatalf("expected 2 products, got %d", len(result.Items))
	}
	assertEqual(t, "Items[0].ID", result.Items[0].ID, "prod-1")
	assertEqual(t, "Items[0].Type", string(result.Items[0].Type), "voucher")
	assertEqual(t, "Items[1].ID", result.Items[1].ID, "prod-2")
	assertEqual(t, "Items[1].Type", string(result.Items[1].Type), "direct_topup")
	assertEqual(t, "HasNext", result.HasNext, false)

	// Verify nested items
	if len(result.Items[0].Items) != 1 {
		t.Fatalf("expected 1 item in product 0, got %d", len(result.Items[0].Items))
	}
	assertEqual(t, "Items[0].Items[0].ID", result.Items[0].Items[0].ID, "item-1")
	assertEqual(t, "Items[0].Items[0].Price", result.Items[0].Items[0].Price, 10.50)
	assertEqual(t, "Items[0].Items[0].Available", result.Items[0].Items[0].Available, true)

	// Verify fields on product 2
	if len(result.Items[1].Fields) != 1 {
		t.Fatalf("expected 1 field in product 1, got %d", len(result.Items[1].Fields))
	}
	assertEqual(t, "Items[1].Fields[0].Key", result.Items[1].Fields[0].Key, "phone")
	assertEqual(t, "Items[1].Fields[0].Required", result.Items[1].Fields[0].Required, true)
}

func TestServicesResource_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "Method", r.Method, "GET")
		assertEqual(t, "Path", r.URL.Path, "/services/prod-1")
		w.Header().Set("Content-Type", "application/json")
		w.Write(successEnvelope(t, newProduct()))
	}))
	defer srv.Close()

	client := approute.NewClient("test-key", approute.WithBaseURL(srv.URL), approute.WithMaxRetries(0))
	product, err := client.Services.Get(t.Context(), "prod-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqual(t, "ID", product.ID, "prod-1")
	if product.Name == nil {
		t.Fatal("Name should not be nil")
	}
	assertEqual(t, "Name", *product.Name, "Steam Wallet 10 USD")
	assertEqual(t, "Type", string(product.Type), "voucher")
}

func TestServicesResource_Stock(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "Method", r.Method, "GET")
		assertEqual(t, "Path", r.URL.Path, "/services/prod-1/stock")
		w.Header().Set("Content-Type", "application/json")
		w.Write(successEnvelope(t, newProductStockResponse(
			newProductStockItem(),
			newProductStockItem(func(s *model.ProductStockItem) {
				s.ItemID = "item-2"
				s.Stock = nil
			}),
		)))
	}))
	defer srv.Close()

	client := approute.NewClient("test-key", approute.WithBaseURL(srv.URL), approute.WithMaxRetries(0))
	stock, err := client.Services.Stock(t.Context(), "prod-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertEqual(t, "ProductID", stock.ProductID, "prod-1")
	if len(stock.Items) != 2 {
		t.Fatalf("expected 2 stock items, got %d", len(stock.Items))
	}
	assertEqual(t, "Items[0].ItemID", stock.Items[0].ItemID, "item-1")
	if stock.Items[0].Stock == nil {
		t.Fatal("Items[0].Stock should not be nil")
	}
	assertEqual(t, "Items[0].Stock", *stock.Items[0].Stock, 100)
	if stock.Items[1].Stock != nil {
		t.Error("Items[1].Stock should be nil")
	}
}

func TestServicesResource_Get_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(errorEnvelope(t, "NOT_FOUND", "Product not found", "t-prod-err-1"))
	}))
	defer srv.Close()

	client := approute.NewClient("test-key", approute.WithBaseURL(srv.URL), approute.WithMaxRetries(0))
	_, err := client.Services.Get(t.Context(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var nf *apierror.NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("expected *NotFoundError, got %T: %v", err, err)
	}
	assertEqual(t, "Code", nf.Code, "NOT_FOUND")
	assertEqual(t, "TraceID", nf.TraceID, "t-prod-err-1")
}
