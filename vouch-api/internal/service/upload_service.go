package service

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// allowedUploadTypes maps accepted MIME types to file extensions.
var allowedUploadTypes = map[string]string{
	"image/png":     "png",
	"image/jpeg":    "jpg",
	"image/webp":    "webp",
	"image/gif":     "gif",
	"image/svg+xml": "svg",
}

// Presigner abstracts object storage presigning (Cloudflare R2 / S3).
type Presigner interface {
	// PresignPut returns a URL the client can PUT to directly, plus the public
	// URL the object will be served from once uploaded.
	PresignPut(ctx context.Context, key, contentType string, ttl time.Duration) (uploadURL, publicURL string, err error)
}

// UploadService issues presigned upload URLs for user assets.
type UploadService struct {
	presigner Presigner
}

// NewUploadService constructs an UploadService.
func NewUploadService(presigner Presigner) *UploadService {
	return &UploadService{presigner: presigner}
}

// PresignResult is returned to the client to perform a direct browser upload.
type PresignResult struct {
	UploadURL string `json:"upload_url"`
	PublicURL string `json:"public_url"`
	Key       string `json:"key"`
	ExpiresIn int    `json:"expires_in"`
}

// PresignUpload validates the content type and returns a presigned PUT URL
// scoped to the given user.
func (s *UploadService) PresignUpload(ctx context.Context, userID, contentType string) (*PresignResult, error) {
	ext, ok := allowedUploadTypes[strings.ToLower(contentType)]
	if !ok {
		return nil, fmt.Errorf("UploadService.PresignUpload: %w", errInvalidUploadType)
	}
	ttl := 5 * time.Minute
	key := fmt.Sprintf("uploads/%s/%s.%s", userID, randHex(12), ext)

	uploadURL, publicURL, err := s.presigner.PresignPut(ctx, key, contentType, ttl)
	if err != nil {
		return nil, fmt.Errorf("UploadService.PresignUpload: %w", err)
	}
	return &PresignResult{
		UploadURL: uploadURL,
		PublicURL: publicURL,
		Key:       key,
		ExpiresIn: int(ttl.Seconds()),
	}, nil
}
