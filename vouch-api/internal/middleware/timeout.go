package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
)

// Timeout returns middleware that cancels the request context after d.
// Handlers that respect c.UserContext() will abort cleanly.
func Timeout(d time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), d)
		defer cancel()
		c.SetUserContext(ctx)

		done := make(chan error, 1)
		go func() { done <- c.Next() }()

		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return response.Error(c, fiber.StatusGatewayTimeout, "timeout", "request timed out")
		}
	}
}
