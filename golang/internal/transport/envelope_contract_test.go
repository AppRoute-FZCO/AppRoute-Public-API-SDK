package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/AppRoute-FZCO/AppRoute-Public-API-SDK/golang/apierror"
	"github.com/AppRoute-FZCO/AppRoute-Public-API-SDK/golang/model"
)

func TestPublicAPIWireEnvelopes(t *testing.T) {
	raw, err := os.ReadFile("../../../testdata/envelopes.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures []struct {
		Name          string             `json:"name"`
		HTTPStatus    int                `json:"httpStatus"`
		Body          json.RawMessage    `json:"body"`
		Data          json.RawMessage    `json:"data"`
		Invalid       bool               `json:"invalid"`
		Code          string             `json:"code"`
		ErrorType     string             `json:"errorType"`
		Message       string             `json:"message"`
		APIStatusCode *int               `json:"apiStatusCode"`
		APIErrorCode  string             `json:"apiErrorCode"`
		Errors        []model.FieldError `json:"errors"`
	}
	if err := json.Unmarshal(raw, &fixtures); err != nil {
		t.Fatal(err)
	}
	for _, fixture := range fixtures {
		t.Run(fixture.Name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(fixture.HTTPStatus)
				w.Write(fixture.Body)
			}))
			defer srv.Close()
			transport := New(srv.URL, "key", srv.Client(), 0)
			data, err := transport.Request(context.Background(), "GET", "/services", nil, nil)
			var base *apierror.ApiError
			if fixture.Invalid {
				if err == nil || !strings.Contains(err.Error(), "envelope") || errors.As(err, &base) {
					t.Fatalf("expected invalid envelope error, got data=%s err=%v", data, err)
				}
				return
			}
			if fixture.Data != nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				var got, want any
				if err := json.Unmarshal(data, &got); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal(fixture.Data, &want); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("data: got %v want %v", got, want)
				}
				return
			}
			if !errors.As(err, &base) {
				t.Fatalf("expected API error, got data=%s err=%v", data, err)
			}
			wantType := fixture.ErrorType + "Error"
			if fixture.ErrorType == "Internal" {
				wantType = "InternalServerError"
			}
			assertEqual(t, "error type", reflect.TypeOf(err).Elem().Name(), wantType)
			assertEqual(t, "code", base.Code, fixture.Code)
			assertEqual(t, "message", base.Message, fixture.Message)
			assertEqual(t, "HTTP status", base.StatusCode, fixture.HTTPStatus)
			var body struct {
				TraceID string `json:"traceId"`
			}
			if err := json.Unmarshal(fixture.Body, &body); err != nil {
				t.Fatal(err)
			}
			assertEqual(t, "trace ID", base.TraceID, body.TraceID)
			if !reflect.DeepEqual(base.APIStatusCode, fixture.APIStatusCode) {
				t.Errorf("API status: got %v want %v", base.APIStatusCode, fixture.APIStatusCode)
			}
			assertEqual(t, "API error code", base.APIErrorCode, fixture.APIErrorCode)
			if !slices.Equal(base.Errors, fixture.Errors) {
				t.Errorf("field errors: got %v want %v", base.Errors, fixture.Errors)
			}
		})
	}
}
