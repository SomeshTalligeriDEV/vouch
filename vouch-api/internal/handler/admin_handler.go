package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
)

// AdminHandler surfaces admin-only endpoints.
type AdminHandler struct {
	svc *service.AdminService
}

// NewAdminHandler constructs an AdminHandler.
func NewAdminHandler(svc *service.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

// Stats handles GET /api/v1/admin/stats.
func (h *AdminHandler) Stats(c *fiber.Ctx) error {
	stats, err := h.svc.Stats(c.UserContext())
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, stats)
}

// ListCompanies handles GET /api/v1/admin/companies.
func (h *AdminHandler) ListCompanies(c *fiber.Ctx) error {
	page := atoiOr(c.Query("page"), 1)
	limit := atoiOr(c.Query("limit"), 20)
	companies, total, err := h.svc.ListCompanies(c.UserContext(), page, limit)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.List(c, companies, page, limit, total)
}
