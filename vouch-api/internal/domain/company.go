package domain

import (
	"context"
	"time"
)

// CompanySize is a rough bucket for headcount.
type CompanySize string

const (
	CompanySizeSolo       CompanySize = "1"
	CompanySizeSmall      CompanySize = "2-10"
	CompanySizeMid        CompanySize = "11-50"
	CompanySizeLarge      CompanySize = "51-200"
	CompanySizeEnterprise CompanySize = "200+"
)

// Company is an organisation that posts problems on the demand board.
// Companies authenticate with email + bcrypt password (not GitHub OAuth).
type Company struct {
	ID           string      `bson:"_id,omitempty" json:"id"`
	Email        string      `bson:"email" json:"email"`
	PasswordHash string      `bson:"password_hash" json:"-"`
	Name         string      `bson:"name" json:"name"`
	Slug         string      `bson:"slug" json:"slug"`
	Website      string      `bson:"website" json:"website"`
	LogoURL      string      `bson:"logo_url" json:"logo_url"`
	Description  string      `bson:"description" json:"description"`
	Size         CompanySize `bson:"size" json:"size"`
	IsVerified   bool        `bson:"is_verified" json:"is_verified"`
	CreatedAt    time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time   `bson:"updated_at" json:"updated_at"`
}

// CompanyRepository abstracts persistence for companies.
type CompanyRepository interface {
	Create(ctx context.Context, c *Company) error
	GetByID(ctx context.Context, id string) (*Company, error)
	GetByEmail(ctx context.Context, email string) (*Company, error)
	GetBySlug(ctx context.Context, slug string) (*Company, error)
	Update(ctx context.Context, c *Company) error
	List(ctx context.Context, page, limit int) ([]*Company, int64, error)
}
