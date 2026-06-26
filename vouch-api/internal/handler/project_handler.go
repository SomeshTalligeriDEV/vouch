package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/middleware"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/validator"
)

// ProjectHandler exposes project endpoints.
type ProjectHandler struct {
	projects *service.ProjectService
	val      *validator.Validator
}

// NewProjectHandler constructs a ProjectHandler.
func NewProjectHandler(projects *service.ProjectService, val *validator.Validator) *ProjectHandler {
	return &ProjectHandler{projects: projects, val: val}
}

// List handles GET /projects.
func (h *ProjectHandler) List(c *fiber.Ctx) error {
	page, limit := pagination(c)
	f := domain.ProjectFilter{
		Status: domain.ProjectStatus(c.Query("status")),
		Tag:    c.Query("tag"),
		Search: c.Query("search"),
		SortBy: c.Query("sort"),
		Page:   page,
		Limit:  limit,
	}
	if v := c.Query("for_sale"); v != "" {
		b := v == "true"
		f.ForSale = &b
	}
	items, total, err := h.projects.List(c.UserContext(), f)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.List(c, items, f.Page, f.Limit, total)
}

type createProjectRequest struct {
	Title       string   `json:"title" validate:"required,min=2,max=120"`
	Tagline     string   `json:"tagline" validate:"max=160"`
	Description string   `json:"description" validate:"max=5000"`
	LogoURL     string   `json:"logo_url" validate:"omitempty,url"`
	LiveURL     string   `json:"live_url" validate:"omitempty,url"`
	RepoURL     string   `json:"repo_url" validate:"omitempty,url"`
	PaymentLink string   `json:"payment_link" validate:"omitempty,url"`
	Tags        []string `json:"tags" validate:"max=10,dive,max=30"`
	ForSale     bool     `json:"for_sale"`
	AskPrice    float64  `json:"ask_price" validate:"gte=0"`
}

// Create handles POST /projects.
func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	var req createProjectRequest
	if err := parseAndValidate(c, h.val, &req); err != nil {
		return err
	}
	p, err := h.projects.Create(c.UserContext(), middleware.UserID(c), service.CreateInput{
		Title:       req.Title,
		Tagline:     req.Tagline,
		Description: req.Description,
		LogoURL:     req.LogoURL,
		LiveURL:     req.LiveURL,
		RepoURL:     req.RepoURL,
		PaymentLink: req.PaymentLink,
		Tags:        req.Tags,
		ForSale:     req.ForSale,
		AskPrice:    req.AskPrice,
	})
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.Created(c, p)
}

// Get handles GET /projects/:id.
func (h *ProjectHandler) Get(c *fiber.Ctx) error {
	p, err := h.projects.Get(c.UserContext(), c.Params("id"))
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, p)
}

type updateProjectRequest struct {
	Title       string   `json:"title" validate:"omitempty,min=2,max=120"`
	Tagline     string   `json:"tagline" validate:"max=160"`
	Description string   `json:"description" validate:"max=5000"`
	LogoURL     string   `json:"logo_url" validate:"omitempty,url"`
	LiveURL     string   `json:"live_url" validate:"omitempty,url"`
	RepoURL     string   `json:"repo_url" validate:"omitempty,url"`
	PaymentLink string   `json:"payment_link" validate:"omitempty,url"`
	Tags        []string `json:"tags" validate:"omitempty,max=10,dive,max=30"`
	ForSale     bool     `json:"for_sale"`
	AskPrice    float64  `json:"ask_price" validate:"gte=0"`
	Status      string   `json:"status" validate:"omitempty,oneof=draft live acquired archived"`
}

// Update handles PATCH /projects/:id.
func (h *ProjectHandler) Update(c *fiber.Ctx) error {
	var req updateProjectRequest
	if err := parseAndValidate(c, h.val, &req); err != nil {
		return err
	}
	p, err := h.projects.Update(c.UserContext(), c.Params("id"), middleware.UserID(c), service.ProjectUpdateInput{
		Title:       req.Title,
		Tagline:     req.Tagline,
		Description: req.Description,
		LogoURL:     req.LogoURL,
		LiveURL:     req.LiveURL,
		RepoURL:     req.RepoURL,
		PaymentLink: req.PaymentLink,
		Tags:        req.Tags,
		ForSale:     req.ForSale,
		AskPrice:    req.AskPrice,
		Status:      domain.ProjectStatus(req.Status),
	})
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, p)
}

// Delete handles DELETE /projects/:id (archive).
func (h *ProjectHandler) Delete(c *fiber.Ctx) error {
	if err := h.projects.Archive(c.UserContext(), c.Params("id"), middleware.UserID(c)); err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, fiber.Map{"archived": true})
}
