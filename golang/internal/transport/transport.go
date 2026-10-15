package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/AppRoute-FZCO/AppRoute-Public-API-SDK/golang/apierror"
	"github.com/AppRoute-FZCO/AppRoute-Public-API-SDK/golang/enum"
	"github.com/AppRoute-FZCO/AppRoute-Public-API-SDK/golang/model"
)

// successCodes are the result codes that indicate a successful API response.
var successCodes = map[enum.ResultCode]struct{}{
	enum.ResultOK:                {},
	enum.ResultAccepted:          {},
	enum.ResultIdempotencyReplay: {},
}

// statusCodes follows the append-only numeric Public API StatusCode contract.
var statusCodes = map[int]enum.ResultCode{
	0:  enum.ResultOK,
	1:  enum.ResultAccepted,
	2:  enum.ResultIdempotencyReplay,
	3:  enum.ResultValidationError,
	4:  enum.ResultUnauthorized,
	5:  enum.ResultForbidden,
	6:  enum.ResultNotFound,
	7:  enum.ResultConflict,
	8:  enum.ResultLimitReached,
	9:  enum.ResultOutOfStock,
	10: enum.ResultInsufficientFunds,
	11: enum.ResultUpstreamError,
	12: enum.ResultInternalError,
	13: enum.ResultAPIKeyLimitExceeded,
	14: enum.ResultProviderUnavailable,
}

// retryableStatusCodes are HTTP status codes eligible for automatic retry.
var retryableStatusCodes = map[int]struct{}{
	429: {},
	500: {},
	502: {},
	503: {},
	504: {},
}

const (
	backoffBase = 1.0 // seconds
)

// HttpTransport handles low-level HTTP communication with the API.
type HttpTransport struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	maxRetries int
}

// New creates a configured HttpTransport.
func New(baseURL, apiKey string, httpClient *http.Client, maxRetries int) *HttpTransport {
	return &HttpTransport{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: httpClient,
		maxRetries: maxRetries,
	}
}

// Request executes an HTTP request against the API, handling retries and
// envelope unwrapping. On success it returns the raw JSON from the "data"
// field; on error it returns a typed *apierror.ApiError (or sub-type).
func (t *HttpTransport) Request(ctx context.Context, method, path string, params url.Values, body any) (json.RawMessage, error) {
	fullURL := t.baseURL + path
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	var bodyReader io.Reader
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("approute: failed to marshal request body: %w", err)
		}
	}

	var lastErr error
	for attempt := 0; attempt <= t.maxRetries; attempt++ {
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("approute: failed to create request: %w", err)
		}

		req.Header.Set("X-API-Key", t.apiKey)
		req.Header.Set("Accept", "application/json")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := t.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("approute: request failed: %w", err)
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("approute: failed to read response body: %w", err)
		}

		// Retry on retryable status codes (except on last attempt).
		if _, ok := retryableStatusCodes[resp.StatusCode]; ok && attempt < t.maxRetries {
			delay := t.retryDelay(resp, attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
				continue
			}
		}

		return t.handleResponse(resp.StatusCode, respBody)
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("approute: max retries exceeded")
}

// handleResponse parses the envelope and either returns data or an error.
func (t *HttpTransport) handleResponse(statusCode int, body []byte) (json.RawMessage, error) {
	if !json.Valid(body) {
		return nil, fmt.Errorf("approute: invalid JSON response (HTTP %d)", statusCode)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil || fields == nil {
		return nil, fmt.Errorf("approute: invalid API response envelope (HTTP %d)", statusCode)
	}
	var env model.Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("approute: invalid API response envelope (HTTP %d): %w", statusCode, err)
	}

	code := enum.ResultCode(env.Code)
	message := env.Message
	if _, present := fields["statusCode"]; present {
		if env.StatusCode == nil {
			return nil, fmt.Errorf("approute: invalid API response envelope (HTTP %d)", statusCode)
		}
		var known bool
		code, known = statusCodes[*env.StatusCode]
		if !known {
			code = enum.ResultCode(strconv.Itoa(*env.StatusCode))
		}
		message = env.StatusMessage
	} else if env.Code == "" {
		return nil, fmt.Errorf("approute: invalid API response envelope (HTTP %d)", statusCode)
	}
	if _, ok := successCodes[code]; ok && statusCode >= 200 && statusCode < 300 && len(env.Errors) == 0 && env.ErrorCode == "" {
		return env.Data, nil
	}

	err := apierror.RaiseForCode(string(code), message, env.TraceID, statusCode, env.Errors)
	var base *apierror.ApiError
	if errors.As(err, &base) {
		base.APIStatusCode = env.StatusCode
		base.APIErrorCode = env.ErrorCode
	}
	return nil, err
}

// retryDelay computes the backoff duration for a given attempt, respecting
// the Retry-After header when present.
func (t *HttpTransport) retryDelay(resp *http.Response, attempt int) time.Duration {
	if ra := resp.Header.Get("Retry-After"); ra != "" {
		if secs, err := strconv.ParseFloat(ra, 64); err == nil {
			return time.Duration(secs * float64(time.Second))
		}
	}
	secs := backoffBase * math.Pow(2, float64(attempt))
	return time.Duration(secs * float64(time.Second))
}
