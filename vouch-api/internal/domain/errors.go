package domain

import "errors"

// Sentinel errors for the domain. Services wrap these with context; handlers
// map them to HTTP status codes.
var (
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidInput      = errors.New("invalid input")
	ErrConflict          = errors.New("conflict")
	ErrProblemClaimed    = errors.New("problem already claimed")
	ErrSelfClaim         = errors.New("cannot claim your own problem")
	ErrStripeNotVerified = errors.New("stripe account not verified")
	ErrInternal          = errors.New("internal error")
)
