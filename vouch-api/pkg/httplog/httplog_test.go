package httplog_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/httplog"
)

func TestFormat_Basic(t *testing.T) {
	e := httplog.Entry{
		Method:     "GET",
		Path:       "/api/v1/problems",
		StatusCode: 200,
		Duration:   123 * time.Millisecond,
	}
	got := e.Format()
	if !strings.Contains(got, "method=GET") {
		t.Errorf("expected method=GET in %q", got)
	}
	if !strings.Contains(got, "status=200") {
		t.Errorf("expected status=200 in %q", got)
	}
	if !strings.Contains(got, "path=/api/v1/problems") {
		t.Errorf("expected path in %q", got)
	}
}

func TestFormat_WithRequestID(t *testing.T) {
	e := httplog.Entry{
		Method:     "POST",
		Path:       "/api/v1/users",
		StatusCode: 201,
		Duration:   5 * time.Millisecond,
		RequestID:  "req-abc123",
	}
	got := e.Format()
	if !strings.Contains(got, "request_id=req-abc123") {
		t.Errorf("expected request_id in %q", got)
	}
}

func TestStatusClass(t *testing.T) {
	cases := map[int]string{
		200: "2xx",
		201: "2xx",
		404: "4xx",
		500: "5xx",
		301: "3xx",
	}
	for code, want := range cases {
		got := httplog.StatusClass(code)
		if got != want {
			t.Errorf("StatusClass(%d) = %q, want %q", code, got, want)
		}
	}
}

func TestIsSuccess(t *testing.T) {
	if !httplog.IsSuccess(200) || !httplog.IsSuccess(204) {
		t.Error("expected 200/204 to be success")
	}
	if httplog.IsSuccess(400) || httplog.IsSuccess(500) {
		t.Error("expected 400/500 to not be success")
	}
}

func TestIsClientError(t *testing.T) {
	if !httplog.IsClientError(400) || !httplog.IsClientError(404) {
		t.Error("expected 400/404 to be client errors")
	}
	if httplog.IsClientError(200) || httplog.IsClientError(500) {
		t.Error("expected 200/500 to not be client errors")
	}
}

func TestIsServerError(t *testing.T) {
	if !httplog.IsServerError(500) || !httplog.IsServerError(503) {
		t.Error("expected 500/503 to be server errors")
	}
	if httplog.IsServerError(200) || httplog.IsServerError(400) {
		t.Error("expected 200/400 to not be server errors")
	}
}

func TestExtractRequestID(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "rid-123")

	got := httplog.ExtractRequestID(req)
	if got != "rid-123" {
		t.Errorf("expected rid-123, got %q", got)
	}
}

func TestExtractRequestID_FallsBackToCorrelation(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Correlation-ID", "corr-456")

	got := httplog.ExtractRequestID(req)
	if got != "corr-456" {
		t.Errorf("expected corr-456, got %q", got)
	}
}

func TestExtractRequestID_Empty(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	got := httplog.ExtractRequestID(req)
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}
