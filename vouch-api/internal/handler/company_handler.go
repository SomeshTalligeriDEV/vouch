package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/middleware"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/validator"
)

// CompanyHandler handles company auth and profile endpoints.
type CompanyHandler struct {
	svc *service.CompanyService
	val *validator.Validator
}

// NewCompanyHandler constructs a CompanyHandler.
func NewCompanyHandler(svc *service.CompanyService, val *validator.Validator) *CompanyHandler {
	return &CompanyHandler{svc: svc, val: val}
}

type companyRegisterInput struct {
	Name     string `json:"name"     validate:"required,min=2,max=120"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Website  string `json:"website"`
	Size     string `json:"size"`
}

// Register handles POST /api/v1/companies/register.
func (h *CompanyHandler) Register(c *fiber.Ctx) error {
	var in companyRegisterInput
	if err := c.BodyParser(&in); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "bad_request", "invalid JSON")
	}
	if err := h.val.Struct(in); err != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation_error", err.Error())
	}
	res, err := h.svc.Register(c.UserContext(), in.Name, in.Email, in.Password, in.Website, in.Size)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.Created(c, res)
}

type companyLoginInput struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Login handles POST /api/v1/companies/login.
func (h *CompanyHandler) Login(c *fiber.Ctx) error {
	var in companyLoginInput
	if err := c.BodyParser(&in); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "bad_request", "invalid JSON")
	}
	if err := h.val.Struct(in); err != nil {
		return response.Error(c, fiber.StatusUnprocessableEntity, "validation_error", err.Error())
	}
	res, err := h.svc.Login(c.UserContext(), in.Email, in.Password)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, res)
}

type companyRefreshInput struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Refresh handles POST /api/v1/companies/refresh.
func (h *CompanyHandler) Refresh(c *fiber.Ctx) error {
	var in companyRefreshInput
	if err := c.BodyParser(&in); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "bad_request", "invalid JSON")
	}
	tokens, err := h.svc.Refresh(c.UserContext(), in.RefreshToken)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, tokens)
}

// GetMe handles GET /api/v1/companies/me — returns the authenticated company.
func (h *CompanyHandler) GetMe(c *fiber.Ctx) error {
	id := middleware.UserID(c)
	company, err := h.svc.GetByID(c.UserContext(), id)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, company)
}

// GetBySlug handles GET /api/v1/companies/:slug — public profile.
func (h *CompanyHandler) GetBySlug(c *fiber.Ctx) error {
	company, err := h.svc.GetBySlug(c.UserContext(), c.Params("slug"))
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, company)
}

type companyUpdateInput struct {
	Name        string `json:"name"`
	Website     string `json:"website"`
	LogoURL     string `json:"logo_url"`
	Description string `json:"description"`
	Size        string `json:"size"`
}

// UpdateMe handles PATCH /api/v1/companies/me.
func (h *CompanyHandler) UpdateMe(c *fiber.Ctx) error {
	var in companyUpdateInput
	if err := c.BodyParser(&in); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "bad_request", "invalid JSON")
	}
	id := middleware.UserID(c)
	company, err := h.svc.UpdateProfile(c.UserContext(), id, in.Name, in.Website, in.LogoURL, in.Description, in.Size)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, company)
}
