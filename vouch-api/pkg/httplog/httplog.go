// Package httplog provides structured HTTP request/response logging helpers.
package httplog

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Entry holds fields captured for one HTTP transaction.
type Entry struct {
	Method     string
	Path       string
	StatusCode int
	Duration   time.Duration
	RemoteAddr string
	UserAgent  string
	RequestID  string
	Bytes      int64
}

// Format returns a structured single-line log string for the entry.
func (e Entry) Format() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("method=%s path=%s status=%d duration=%s",
		e.Method, e.Path, e.StatusCode, e.Duration.Round(time.Millisecond)))
	if e.RemoteAddr != "" {
		sb.WriteString(fmt.Sprintf(" remote=%s", e.RemoteAddr))
	}
	if e.RequestID != "" {
		sb.WriteString(fmt.Sprintf(" request_id=%s", e.RequestID))
	}
	if e.Bytes > 0 {
		sb.WriteString(fmt.Sprintf(" bytes=%d", e.Bytes))
	}
	return sb.String()
}

// StatusClass returns "2xx", "4xx", "5xx", etc. for a status code.
func StatusClass(code int) string {
	return fmt.Sprintf("%dxx", code/100)
}

// IsSuccess returns true for 2xx status codes.
func IsSuccess(code int) bool {
	return code >= 200 && code < 300
}

// IsClientError returns true for 4xx status codes.
func IsClientError(code int) bool {
	return code >= 400 && code < 500
}

// IsServerError returns true for 5xx status codes.
func IsServerError(code int) bool {
	return code >= 500 && code < 600
}

// ExtractRequestID reads the X-Request-ID header, falling back to X-Correlation-ID.
func ExtractRequestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	return r.Header.Get("X-Correlation-ID")
}
