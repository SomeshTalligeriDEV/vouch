package observability

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// RequestLogger is a minimal structured request logger for Fiber.
// For production, wire up Prometheus or Datadog here.
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		log.Info().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Str("ip", c.IP()).
			Str("request_id", c.Locals("requestID").(string)).
			Msg("request")

		return err
	}
}
