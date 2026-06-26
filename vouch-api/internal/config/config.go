package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all runtime configuration, loaded from the environment.
type Config struct {
	MongoURI    string
	MongoDB     string
	RedisURL    string

	JWTSecret        string
	JWTRefreshSecret string

	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string

	StripeClientID  string
	StripeSecretKey string

	ResendAPIKey string
	EmailFrom    string

	R2Bucket    string
	R2AccessKey string
	R2SecretKey string
	R2Endpoint  string
	R2PublicURL string

	SentryDSN string

	AppURL string

	Port string
	Env  string
}

// IsProduction reports whether the service is running in production.
func (c *Config) IsProduction() bool { return c.Env == "production" }

// Load reads configuration from environment variables and validates that all
// required values are present.
func Load() (*Config, error) {
	c := &Config{
		MongoURI:           os.Getenv("MONGO_URI"),
		MongoDB:            getOr("MONGO_DB", "vouch"),
		RedisURL:           os.Getenv("REDIS_URL"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTRefreshSecret:   os.Getenv("JWT_REFRESH_SECRET"),
		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GitHubRedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		StripeClientID:     os.Getenv("STRIPE_CLIENT_ID"),
		StripeSecretKey:    os.Getenv("STRIPE_SECRET_KEY"),
		ResendAPIKey:       os.Getenv("RESEND_API_KEY"),
		EmailFrom:          getOr("EMAIL_FROM", "Vouch <noreply@vouch.dev>"),
		R2Bucket:           os.Getenv("R2_BUCKET"),
		R2AccessKey:        os.Getenv("R2_ACCESS_KEY"),
		R2SecretKey:        os.Getenv("R2_SECRET_KEY"),
		R2Endpoint:         os.Getenv("R2_ENDPOINT"),
		R2PublicURL:        os.Getenv("R2_PUBLIC_URL"),
		SentryDSN:          os.Getenv("SENTRY_DSN"),
		AppURL:             getOr("APP_URL", "http://localhost:3000"),
		Port:               getOr("PORT", "8080"),
		Env:                getOr("ENV", "development"),
	}

	required := map[string]string{
		"MONGO_URI":          c.MongoURI,
		"REDIS_URL":          c.RedisURL,
		"JWT_SECRET":         c.JWTSecret,
		"JWT_REFRESH_SECRET": c.JWTRefreshSecret,
	}
	var missing []string
	for k, v := range required {
		if strings.TrimSpace(v) == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("config.Load: missing required env vars: %s", strings.Join(missing, ", "))
	}
	return c, nil
}

func getOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
