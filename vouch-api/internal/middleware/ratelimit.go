package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
)

// RateLimiter is a fixed-window Redis-backed limiter.
type RateLimiter struct {
	rdb    *redis.Client
	limit  int
	window time.Duration
}

// NewRateLimiter constructs a limiter allowing `limit` requests per `window`.
func NewRateLimiter(rdb *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{rdb: rdb, limit: limit, window: window}
}

// Limit returns middleware enforcing the rate limit, keyed by authenticated
// user id when present, otherwise by client IP.
// It also sets X-RateLimit-Limit and X-RateLimit-Remaining response headers
// so clients can back off gracefully.
func (rl *RateLimiter) Limit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := UserID(c)
		if id == "" {
			id = c.IP()
		}
		key := fmt.Sprintf("ratelimit:%s:%s", c.Path(), id)

		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		count, err := rl.rdb.Incr(ctx, key).Result()
		if err != nil {
			// Fail open: never block traffic because Redis hiccupped.
			return c.Next()
		}
		if count == 1 {
			rl.rdb.Expire(ctx, key, rl.window) //nolint:errcheck
		}

		remaining := int64(rl.limit) - count
		if remaining < 0 {
			remaining = 0
		}
		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(rl.window).Unix()))

		if count > int64(rl.limit) {
			c.Set("Retry-After", fmt.Sprintf("%.0f", rl.window.Seconds()))
			return response.Error(c, fiber.StatusTooManyRequests, "rate_limited", "too many requests")
		}
		return c.Next()
	}
}
