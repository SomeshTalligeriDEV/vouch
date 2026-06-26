package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"
)

const headerRequestID = "X-Request-ID"

// RequestID attaches a unique request ID to every request for log correlation.
// It honours an incoming X-Request-ID header (from a load balancer or client)
// and always echoes the final value back in the response.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get(headerRequestID)
		if id == "" {
			id = newRequestID()
		}
		c.Locals("requestID", id)
		c.Set(headerRequestID, id)
		return c.Next()
	}
}

func newRequestID() string {
	b := make([]byte, 8)
	rand.Read(b) //nolint:errcheck
	return hex.EncodeToString(b)
}
