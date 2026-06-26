package response

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// Envelope is the standard JSON shape for all API responses.
type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorBody describes a failure in a stable, machine-readable way.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta carries pagination metadata for list endpoints.
type Meta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// OK writes a 200 response with data.
func OK(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Envelope{Success: true, Data: data})
}

// Created writes a 201 response with data.
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Envelope{Success: true, Data: data})
}

// List writes a 200 response with data and pagination metadata.
func List(c *fiber.Ctx, data interface{}, page, limit int, total int64) error {
	return c.Status(fiber.StatusOK).JSON(Envelope{
		Success: true,
		Data:    data,
		Meta:    &Meta{Page: page, Limit: limit, Total: total},
	})
}

// Error writes a failure response with an explicit status and code.
func Error(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(Envelope{
		Success: false,
		Error:   &ErrorBody{Code: code, Message: message},
	})
}

// FromDomain maps a domain error to an appropriate HTTP response.
func FromDomain(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return Error(c, fiber.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, domain.ErrAlreadyExists), errors.Is(err, domain.ErrConflict),
		errors.Is(err, domain.ErrProblemClaimed):
		return Error(c, fiber.StatusConflict, "conflict", err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		return Error(c, fiber.StatusUnauthorized, "unauthorized", err.Error())
	case errors.Is(err, domain.ErrForbidden), errors.Is(err, domain.ErrSelfClaim):
		return Error(c, fiber.StatusForbidden, "forbidden", err.Error())
	case errors.Is(err, domain.ErrInvalidInput), errors.Is(err, domain.ErrStripeNotVerified):
		return Error(c, fiber.StatusBadRequest, "invalid_input", err.Error())
	default:
		return Error(c, fiber.StatusInternalServerError, "internal_error", "something went wrong")
	}
}
