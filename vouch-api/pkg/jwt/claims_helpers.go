package jwt

import "time"

// IsExpired reports whether the token's expiry time has passed.
func (c *Claims) IsExpired() bool {
	if c.ExpiresAt == nil {
		return true
	}
	return time.Now().After(c.ExpiresAt.Time)
}

// IsCompany reports whether the token was issued for a company account.
func (c *Claims) IsCompany() bool { return c.SubjectType == SubjectTypeCompany }

// IsUser reports whether the token was issued for a builder account.
func (c *Claims) IsUser() bool { return c.SubjectType == SubjectTypeUser || c.SubjectType == "" }
