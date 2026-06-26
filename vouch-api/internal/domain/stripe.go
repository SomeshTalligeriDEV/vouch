package domain

import "time"

// StripeSnapshot is a point-in-time read of a builder's Stripe revenue data.
// Vouch reads this via read-only OAuth; it never processes payments.
type StripeSnapshot struct {
	ID             string    `bson:"_id,omitempty" json:"id"`
	BuilderID      string    `bson:"builder_id" json:"builder_id"`
	MRR            float64   `bson:"mrr" json:"mrr"`
	TotalCustomers int       `bson:"total_customers" json:"total_customers"`
	Currency       string    `bson:"currency" json:"currency"`
	VerifiedAt     time.Time `bson:"verified_at" json:"verified_at"`
}
