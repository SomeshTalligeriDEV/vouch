package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/middleware"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/validator"
)

// ReviewHandler exposes review endpoints.
type ReviewHandler struct {
	reviews *service.ReviewService
	val     *validator.Validator
}

// NewReviewHandler constructs a ReviewHandler.
func NewReviewHandler(reviews *service.ReviewService, val *validator.Validator) *ReviewHandler {
	return &ReviewHandler{reviews: reviews, val: val}
}

type createReviewRequest struct {
	ProjectID string `json:"project_id" validate:"required"`
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	Body      string `json:"body" validate:"max=2000"`
}

// Create handles POST /reviews.
func (h *ReviewHandler) Create(c *fiber.Ctx) error {
	var req createReviewRequest
	if err := parseAndValidate(c, h.val, &req); err != nil {
		return err
	}
	r, err := h.reviews.Create(c.UserContext(), middleware.UserID(c), service.ReviewInput{
		ProjectID: req.ProjectID,
		Rating:    req.Rating,
		Body:      req.Body,
	})
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.Created(c, r)
}

// ListByProject handles GET /reviews/project/:id.
func (h *ReviewHandler) ListByProject(c *fiber.Ctx) error {
	page, limit := pagination(c)
	items, total, err := h.reviews.ListByProject(c.UserContext(), c.Params("id"), page, limit)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.List(c, items, page, limit, total)
}
