package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/middleware"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/validator"
)

// ProblemHandler exposes demand-board endpoints.
type ProblemHandler struct {
	problems *service.ProblemService
	val      *validator.Validator
}

// NewProblemHandler constructs a ProblemHandler.
func NewProblemHandler(problems *service.ProblemService, val *validator.Validator) *ProblemHandler {
	return &ProblemHandler{problems: problems, val: val}
}

// List handles GET /problems.
func (h *ProblemHandler) List(c *fiber.Ctx) error {
	page, limit := pagination(c)
	f := domain.ProblemFilter{
		Status: domain.ProblemStatus(c.Query("status")),
		Tag:    c.Query("tag"),
		Search: c.Query("search"),
		SortBy: c.Query("sort"),
		Page:   page,
		Limit:  limit,
	}
	items, total, err := h.problems.List(c.UserContext(), f)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.List(c, items, f.Page, f.Limit, total)
}

type createProblemRequest struct {
	Title       string   `json:"title" validate:"required,min=4,max=140"`
	Description string   `json:"description" validate:"required,max=5000"`
	Tags        []string `json:"tags" validate:"max=10,dive,max=30"`
	BudgetMin   float64  `json:"budget_min" validate:"gte=0"`
	BudgetMax   float64  `json:"budget_max" validate:"gte=0"`
}

// Create handles POST /problems.
func (h *ProblemHandler) Create(c *fiber.Ctx) error {
	var req createProblemRequest
	if err := parseAndValidate(c, h.val, &req); err != nil {
		return err
	}
	p, err := h.problems.Create(c.UserContext(), middleware.UserID(c), service.ProblemInput{
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		BudgetMin:   req.BudgetMin,
		BudgetMax:   req.BudgetMax,
	})
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.Created(c, p)
}

// Get handles GET /problems/:id.
func (h *ProblemHandler) Get(c *fiber.Ctx) error {
	p, err := h.problems.Get(c.UserContext(), c.Params("id"))
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, p)
}

// Claim handles POST /problems/:id/claim.
func (h *ProblemHandler) Claim(c *fiber.Ctx) error {
	p, err := h.problems.Claim(c.UserContext(), c.Params("id"), middleware.UserID(c))
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, p)
}

// Upvote handles POST /problems/:id/upvote.
func (h *ProblemHandler) Upvote(c *fiber.Ctx) error {
	p, err := h.problems.Upvote(c.UserContext(), c.Params("id"), middleware.UserID(c))
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, p)
}
