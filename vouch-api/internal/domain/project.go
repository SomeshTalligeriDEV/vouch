package domain

import "time"

// ProjectStatus enumerates the lifecycle states of a project.
type ProjectStatus string

const (
	ProjectStatusDraft    ProjectStatus = "draft"
	ProjectStatusLive     ProjectStatus = "live"
	ProjectStatusAcquired ProjectStatus = "acquired"
	ProjectStatusArchived ProjectStatus = "archived"
)

// Project is a shipped (or for-sale) product owned by a builder.
type Project struct {
	ID            string        `bson:"_id,omitempty" json:"id"`
	BuilderID     string        `bson:"builder_id" json:"builder_id"`
	Title         string        `bson:"title" json:"title"`
	Slug          string        `bson:"slug" json:"slug"`
	Tagline       string        `bson:"tagline" json:"tagline"`
	Description    string       `bson:"description" json:"description"`
	LogoURL       string        `bson:"logo_url" json:"logo_url"`
	LiveURL       string        `bson:"live_url" json:"live_url"`
	RepoURL       string        `bson:"repo_url" json:"repo_url"`
	PaymentLink   string        `bson:"payment_link" json:"payment_link"`
	Tags          []string      `bson:"tags" json:"tags"`
	Status        ProjectStatus `bson:"status" json:"status"`
	ForSale       bool          `bson:"for_sale" json:"for_sale"`
	AskPrice      float64       `bson:"ask_price" json:"ask_price"`
	VerifiedUsers int           `bson:"verified_users" json:"verified_users"`
	MRR           float64       `bson:"mrr" json:"mrr"`
	ReviewCount   int           `bson:"review_count" json:"review_count"`
	AverageRating float64       `bson:"average_rating" json:"average_rating"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at"`
}

// ProjectFilter holds optional query constraints for listing projects.
type ProjectFilter struct {
	BuilderID string
	Status    ProjectStatus
	ForSale   *bool
	Tag       string
	Search    string
	Page      int
	Limit     int
	SortBy    string // "mrr", "users", "rating", "recent"
}
