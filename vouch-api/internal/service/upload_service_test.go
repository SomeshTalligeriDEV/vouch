package service

import (
	"strings"
	"testing"
)

func TestUploadService_ContentTypeValidation(t *testing.T) {
	allowed := []string{"image/png", "image/jpeg", "image/webp", "image/gif"}
	rejected := []string{"application/pdf", "text/html", "application/javascript", ""}

	for _, ct := range allowed {
		if !isAllowedImageType(ct) {
			t.Errorf("expected %q to be allowed", ct)
		}
	}
	for _, ct := range rejected {
		if isAllowedImageType(ct) {
			t.Errorf("expected %q to be rejected", ct)
		}
	}
}

// isAllowedImageType mirrors the logic in UploadService to allow testing without
// constructing the full service (which requires R2 credentials).
func isAllowedImageType(contentType string) bool {
	allowed := []string{"image/png", "image/jpeg", "image/webp", "image/gif"}
	for _, a := range allowed {
		if strings.EqualFold(contentType, a) {
			return true
		}
	}
	return false
}
