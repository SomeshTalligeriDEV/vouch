package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/validator"
)

// parseAndValidate decodes the JSON body into dst and validates it, writing a
// 400 response and returning a non-nil error on failure.
func parseAndValidate(c *fiber.Ctx, val *validator.Validator, dst interface{}) error {
	if err := c.BodyParser(dst); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid_body", "malformed request body")
	}
	if err := val.Struct(dst); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "validation_error", err.Error())
	}
	return nil
}

// pagination reads page and limit query params with sane defaults.
func pagination(c *fiber.Ctx) (page, limit int) {
	page = atoiOr(c.Query("page"), 1)
	limit = atoiOr(c.Query("limit"), 20)
	return page, limit
}

func atoiOr(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}
