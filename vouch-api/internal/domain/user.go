package domain

import (
	"strings"
	"time"
)

// Role enumerates the kinds of accounts on Vouch.
type Role string

const (
	RoleBuilder Role = "builder"
	RoleUser    Role = "user"
	RoleAdmin   Role = "admin"
)

// User is a person on the platform. A user can both post problems and build.
type User struct {
	ID              string    `bson:"_id,omitempty" json:"id"`
	Email           string    `bson:"email" json:"email"`
	Username        string    `bson:"username" json:"username"`
	Name            string    `bson:"name" json:"name"`
	Bio             string    `bson:"bio" json:"bio"`
	AvatarURL       string    `bson:"avatar_url" json:"avatar_url"`
	GitHubID        int64     `bson:"github_id" json:"github_id"`
	GitHubLogin     string    `bson:"github_login" json:"github_login"`
	StripeAccountID string    `bson:"stripe_account_id" json:"-"` // hidden from API responses
	Role            Role      `bson:"role" json:"role"`
	IsVerified      bool      `bson:"is_verified" json:"is_verified"`
	WebsiteURL      string    `bson:"website_url" json:"website_url"`
	TwitterHandle   string    `bson:"twitter_handle" json:"twitter_handle"`
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time `bson:"updated_at" json:"updated_at"`
}

// DisplayName returns the user's full name when set, otherwise their username.
func (u *User) DisplayName() string {
	if strings.TrimSpace(u.Name) != "" {
		return u.Name
	}
	return u.Username
}

// IsAdmin reports whether the user holds the admin role.
func (u *User) IsAdmin() bool { return u.Role == RoleAdmin }

// HasStripe reports whether the user has connected a Stripe account.
func (u *User) HasStripe() bool {
	return u.StripeAccountID != ""
}
