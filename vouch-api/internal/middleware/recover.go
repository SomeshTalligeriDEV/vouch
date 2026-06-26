package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
)

// Recover returns middleware that catches panics, logs them with a full stack
// trace, and returns a 500 to the client instead of crashing the process.
func Recover(log zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				log.Error().
					Str("panic", fmt.Sprintf("%v", r)).
					Str("stack", string(stack)).
					Str("path", c.Path()).
					Str("method", c.Method()).
					Msg("recovered from panic")
				err = response.Error(c, fiber.StatusInternalServerError, "internal_error", "something went wrong")
			}
		}()
		return c.Next()
	}
}
