package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// StripeClient implements service.StripeGateway against the real Stripe API
// using read-only Connect OAuth. Vouch never creates charges.
type StripeClient struct {
	clientID  string
	secretKey string
	http      *http.Client
}

// NewStripeClient constructs a StripeClient.
func NewStripeClient(clientID, secretKey string) *StripeClient {
	return &StripeClient{
		clientID:  clientID,
		secretKey: secretKey,
		http:      &http.Client{Timeout: 15 * time.Second},
	}
}

// ExchangeCode swaps an OAuth authorization code for a connected account id.
func (s *StripeClient) ExchangeCode(ctx context.Context, code string) (string, error) {
	form := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://connect.stripe.com/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("StripeClient.ExchangeCode: %w", err)
	}
	req.SetBasicAuth(s.secretKey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("StripeClient.ExchangeCode: %w", err)
	}
	defer resp.Body.Close()

	var body struct {
		StripeUserID     string `json:"stripe_user_id"`
		ErrorDescription string `json:"error_description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("StripeClient.ExchangeCode: %w", err)
	}
	if body.StripeUserID == "" {
		return "", fmt.Errorf("StripeClient.ExchangeCode: %s", body.ErrorDescription)
	}
	return body.StripeUserID, nil
}

// FetchRevenue reads active subscription MRR and customer count for an account.
func (s *StripeClient) FetchRevenue(ctx context.Context, accountID string) (float64, int, string, error) {
	// List active subscriptions on behalf of the connected account.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.stripe.com/v1/subscriptions?status=active&limit=100", nil)
	if err != nil {
		return 0, 0, "", fmt.Errorf("StripeClient.FetchRevenue: %w", err)
	}
	req.SetBasicAuth(s.secretKey, "")
	req.Header.Set("Stripe-Account", accountID)

	resp, err := s.http.Do(req)
	if err != nil {
		return 0, 0, "", fmt.Errorf("StripeClient.FetchRevenue: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, 0, "", fmt.Errorf("StripeClient.FetchRevenue: stripe status %d", resp.StatusCode)
	}

	var body struct {
		Data []struct {
			Currency string `json:"currency"`
			Items    struct {
				Data []struct {
					Price struct {
						UnitAmount int64  `json:"unit_amount"`
						Recurring  struct {
							Interval      string `json:"interval"`
							IntervalCount int    `json:"interval_count"`
						} `json:"recurring"`
					} `json:"price"`
					Quantity int64 `json:"quantity"`
				} `json:"data"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, 0, "", fmt.Errorf("StripeClient.FetchRevenue: %w", err)
	}

	var mrrCents float64
	currency := "usd"
	for _, sub := range body.Data {
		if sub.Currency != "" {
			currency = sub.Currency
		}
		for _, it := range sub.Items.Data {
			qty := it.Quantity
			if qty == 0 {
				qty = 1
			}
			mrrCents += normalizeToMonthly(
				float64(it.Price.UnitAmount*qty),
				it.Price.Recurring.Interval,
				it.Price.Recurring.IntervalCount,
			)
		}
	}
	return mrrCents / 100.0, len(body.Data), currency, nil
}

// normalizeToMonthly converts a recurring amount to its monthly equivalent.
func normalizeToMonthly(amount float64, interval string, count int) float64 {
	if count <= 0 {
		count = 1
	}
	switch interval {
	case "year":
		return amount / (12.0 * float64(count))
	case "week":
		return amount * (52.0 / 12.0) / float64(count)
	case "day":
		return amount * (365.0 / 12.0) / float64(count)
	default: // month
		return amount / float64(count)
	}
}
