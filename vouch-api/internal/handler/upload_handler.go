package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/middleware"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/validator"
)

// UploadHandler exposes the presigned-upload endpoint.
type UploadHandler struct {
	uploads *service.UploadService
	val     *validator.Validator
}

// NewUploadHandler constructs an UploadHandler.
func NewUploadHandler(uploads *service.UploadService, val *validator.Validator) *UploadHandler {
	return &UploadHandler{uploads: uploads, val: val}
}

type presignRequest struct {
	ContentType string `json:"content_type" validate:"required"`
}

// Presign handles POST /uploads/presign.
func (h *UploadHandler) Presign(c *fiber.Ctx) error {
	var req presignRequest
	if err := parseAndValidate(c, h.val, &req); err != nil {
		return err
	}
	res, err := h.uploads.PresignUpload(c.UserContext(), middleware.UserID(c), req.ContentType)
	if err != nil {
		return response.FromDomain(c, err)
	}
	return response.OK(c, res)
}
