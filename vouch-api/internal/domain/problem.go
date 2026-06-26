package domain

import "time"

// ProblemStatus enumerates the lifecycle of a demand-board problem.
type ProblemStatus string

const (
	ProblemStatusOpen      ProblemStatus = "open"
	ProblemStatusClaimed   ProblemStatus = "claimed"
	ProblemStatusShipped   ProblemStatus = "shipped"
	ProblemStatusCancelled ProblemStatus = "cancelled"
)

// Problem is a real demand posted by a user, with a budget, that a builder
// can claim and ship.
type Problem struct {
	ID          string        `bson:"_id,omitempty" json:"id"`
	PosterID    string        `bson:"poster_id" json:"poster_id"`
	ClaimedBy   string        `bson:"claimed_by,omitempty" json:"claimed_by,omitempty"`
	ShippedProjectID string   `bson:"shipped_project_id,omitempty" json:"shipped_project_id,omitempty"`
	Title       string        `bson:"title" json:"title"`
	Slug        string        `bson:"slug" json:"slug"`
	Description string        `bson:"description" json:"description"`
	Tags        []string      `bson:"tags" json:"tags"`
	BudgetMin   float64       `bson:"budget_min" json:"budget_min"`
	BudgetMax   float64       `bson:"budget_max" json:"budget_max"`
	Status      ProblemStatus `bson:"status" json:"status"`
	Upvotes     int           `bson:"upvotes" json:"upvotes"`
	UpvotedBy   []string      `bson:"upvoted_by" json:"-"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
}

// ProblemFilter holds optional query constraints for listing problems.
type ProblemFilter struct {
	PosterID  string
	ClaimedBy string
	Status    ProblemStatus
	Tag       string
	Search    string
	Page      int
	Limit     int
	SortBy    string // "upvotes", "budget", "recent"
}

// HasUpvoted reports whether the given user already upvoted this problem.
func (p *Problem) HasUpvoted(userID string) bool {
	for _, id := range p.UpvotedBy {
		if id == userID {
			return true
		}
	}
	return false
}
