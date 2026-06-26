package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/middleware"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
)

// ScoreHandler exposes builder score endpoints.
type ScoreHandler struct {
	scores *service.ScoreService
}

// NewScoreHandler constructs a ScoreHandler.
func NewScoreHandler(scores *service.ScoreService) *ScoreHandler {
	return &ScoreHandler{scores: scores}
}

// GetByUsername handles GET /scores/:username.
func (h *ScoreHandler) GetByUsername(c *fiber.Ctx) error {
	score, err := h.scores.GetByUsername(c.UserContext(), c.Params("username"))
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, score)
}

// Recalculate handles POST /scores/recalculate for the authenticated builder.
func (h *ScoreHandler) Recalculate(c *fiber.Ctx) error {
	if err := h.scores.EnqueueRecalc(c.UserContext(), middleware.UserID(c)); err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, fiber.Map{"queued": true})
}

// Leaderboard handles GET /scores (top builders).
func (h *ScoreHandler) Leaderboard(c *fiber.Ctx) error {
	out, err := h.scores.Leaderboard(c.UserContext(), atoiOr(c.Query("limit"), 25))
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, out)
}
