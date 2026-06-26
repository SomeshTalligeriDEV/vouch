package domain

import "time"

// Review is a rating left by a user on a project.
type Review struct {
	ID               string    `bson:"_id,omitempty" json:"id"`
	ProjectID        string    `bson:"project_id" json:"project_id"`
	ReviewerID       string    `bson:"reviewer_id" json:"reviewer_id"`
	ReviewerUsername string    `bson:"reviewer_username" json:"reviewer_username"`
	Rating           int       `bson:"rating" json:"rating"` // 1..5
	Body             string    `bson:"body" json:"body"`
	VerifiedPurchase bool      `bson:"verified_purchase" json:"verified_purchase"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
}

// ReviewStats is an aggregate over a project's reviews.
type ReviewStats struct {
	Count   int
	Average float64
}

// IsVerifiedPurchase reports whether the reviewer paid for the project via Stripe.
func (r *Review) IsVerifiedPurchase() bool { return r.VerifiedPurchase }

// IsValid reports whether the review has a rating in the valid 1–5 range.
func (r *Review) IsValid() bool { return r.Rating >= 1 && r.Rating <= 5 }
