package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/middleware"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/validator"
)

// UserHandler exposes user and auth endpoints.
type UserHandler struct {
	users  *service.UserService
	stripe *service.StripeService
	val    *validator.Validator
}

// NewUserHandler constructs a UserHandler.
func NewUserHandler(users *service.UserService, stripe *service.StripeService, val *validator.Validator) *UserHandler {
	return &UserHandler{users: users, stripe: stripe, val: val}
}

type githubAuthRequest struct {
	Code string `json:"code" validate:"required"`
}

type authResponse struct {
	User   interface{} `json:"user"`
	Tokens interface{} `json:"tokens"`
}

// GitHubCallback handles POST /auth/github.
func (h *UserHandler) GitHubCallback(c *fiber.Ctx) error {
	var req githubAuthRequest
	if err := parseAndValidate(c, h.val, &req); err != nil {
		return err
	}
	user, tokens, err := h.users.LoginWithGitHub(c.UserContext(), req.Code)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, authResponse{User: user, Tokens: tokens})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Refresh handles POST /auth/refresh.
func (h *UserHandler) Refresh(c *fiber.Ctx) error {
	var req refreshRequest
	if err := parseAndValidate(c, h.val, &req); err != nil {
		return err
	}
	tokens, err := h.users.Refresh(c.UserContext(), req.RefreshToken)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, tokens)
}

// GetByUsername handles GET /users/:username.
func (h *UserHandler) GetByUsername(c *fiber.Ctx) error {
	user, err := h.users.GetByUsername(c.UserContext(), c.Params("username"))
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, user)
}

// GetMe handles GET /users/me — returns the authenticated user's own profile.
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	username := middleware.Username(c)
	if username == "" {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized", "missing identity")
	}
	user, err := h.users.GetByUsername(c.UserContext(), username)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, user)
}

type updateProfileRequest struct {
	Name          string `json:"name" validate:"max=100"`
	Bio           string `json:"bio" validate:"max=500"`
	AvatarURL     string `json:"avatar_url" validate:"omitempty,url"`
	WebsiteURL    string `json:"website_url" validate:"omitempty,url"`
	TwitterHandle string `json:"twitter_handle" validate:"max=50"`
}

// UpdateMe handles PATCH /users/me.
func (h *UserHandler) UpdateMe(c *fiber.Ctx) error {
	var req updateProfileRequest
	if err := parseAndValidate(c, h.val, &req); err != nil {
		return err
	}
	user, err := h.users.UpdateProfile(c.UserContext(), middleware.UserID(c), service.UpdateInput{
		Name:          req.Name,
		Bio:           req.Bio,
		AvatarURL:     req.AvatarURL,
		WebsiteURL:    req.WebsiteURL,
		TwitterHandle: req.TwitterHandle,
	})
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, user)
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Logout handles POST /auth/logout — revokes the supplied refresh token.
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	var req logoutRequest
	if err := parseAndValidate(c, h.val, &req); err != nil {
		return err
	}
	if err := h.users.Logout(c.UserContext(), req.RefreshToken); err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, fiber.Map{"logged_out": true})
}

type connectStripeRequest struct {
	Code string `json:"code" validate:"required"`
}

// ConnectStripe handles POST /users/me/stripe.
func (h *UserHandler) ConnectStripe(c *fiber.Ctx) error {
	var req connectStripeRequest
	if err := parseAndValidate(c, h.val, &req); err != nil {
		return err
	}
	snap, err := h.stripe.Connect(c.UserContext(), middleware.UserID(c), req.Code)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, snap)
}
