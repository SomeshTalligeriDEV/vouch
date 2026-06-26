package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ResendClient implements service.Notifier against the Resend REST API.
type ResendClient struct {
	apiKey string
	from   string
	http   *http.Client
}

// NewResendClient constructs a ResendClient. `from` is the verified sender,
// e.g. "Vouch <noreply@vouch.dev>".
func NewResendClient(apiKey, from string) *ResendClient {
	return &ResendClient{
		apiKey: apiKey,
		from:   from,
		http:   &http.Client{Timeout: 10 * time.Second},
	}
}

// Send delivers a single HTML email.
func (r *ResendClient) Send(ctx context.Context, to, subject, htmlBody string) error {
	if r.apiKey == "" {
		// No key configured (e.g. local dev) — treat as a no-op success.
		return nil
	}
	payload, err := json.Marshal(map[string]any{
		"from":    r.from,
		"to":      []string{to},
		"subject": subject,
		"html":    htmlBody,
	})
	if err != nil {
		return fmt.Errorf("ResendClient.Send: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.resend.com/emails", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("ResendClient.Send: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.http.Do(req)
	if err != nil {
		return fmt.Errorf("ResendClient.Send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("ResendClient.Send: resend status %d", resp.StatusCode)
	}
	return nil
}
