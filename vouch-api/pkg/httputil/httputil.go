// Package httputil provides HTTP helper utilities for Vouch.
package httputil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultClient is a pre-configured HTTP client with reasonable timeouts.
var DefaultClient = &http.Client{
	Timeout: 10 * time.Second,
}

// GetJSON makes a GET request and decodes the JSON response into dst.
func GetJSON(ctx context.Context, url string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("httputil.GetJSON: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("httputil.GetJSON: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("httputil.GetJSON: status %d: %s", resp.StatusCode, body)
	}

	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		return fmt.Errorf("httputil.GetJSON: decode: %w", err)
	}
	return nil
}

// PostJSON makes a POST request with a JSON body and decodes the response.
func PostJSON(ctx context.Context, url string, body any, dst any) error {
	pr, pw := io.Pipe()
	go func() {
		enc := json.NewEncoder(pw)
		_ = pw.CloseWithError(enc.Encode(body))
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, pr)
	if err != nil {
		return fmt.Errorf("httputil.PostJSON: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("httputil.PostJSON: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("httputil.PostJSON: status %d: %s", resp.StatusCode, errBody)
	}

	if dst != nil {
		if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
			return fmt.Errorf("httputil.PostJSON: decode: %w", err)
		}
	}
	return nil
}

// BearerHeader returns an Authorization header value for the given token.
func BearerHeader(token string) string {
	return "Bearer " + token
}
