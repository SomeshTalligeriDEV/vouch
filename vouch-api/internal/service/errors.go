package service

import (
	"errors"
	"fmt"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// errInvalidUploadType wraps domain.ErrInvalidInput for rejected MIME types so
// handlers map it to a 400.
var errInvalidUploadType = fmt.Errorf("unsupported file type: %w", domain.ErrInvalidInput)

// isNotFound reports whether err wraps domain.ErrNotFound.
func isNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}
