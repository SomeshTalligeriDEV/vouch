package httputil_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/pkg/httputil"
)

func TestGetJSON_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"hello": "world"})
	}))
	defer srv.Close()

	var result map[string]string
	if err := httputil.GetJSON(context.Background(), srv.URL, &result); err != nil {
		t.Fatalf("GetJSON: %v", err)
	}
	if result["hello"] != "world" {
		t.Errorf("expected hello=world, got %q", result["hello"])
	}
}

func TestGetJSON_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer srv.Close()

	var dst any
	err := httputil.GetJSON(context.Background(), srv.URL, &dst)
	if err == nil {
		t.Error("expected error for 404 response")
	}
}

func TestPostJSON_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"echo": body["input"]})
	}))
	defer srv.Close()

	var result map[string]string
	if err := httputil.PostJSON(context.Background(), srv.URL, map[string]string{"input": "ping"}, &result); err != nil {
		t.Fatalf("PostJSON: %v", err)
	}
	if result["echo"] != "ping" {
		t.Errorf("expected echo=ping, got %q", result["echo"])
	}
}

func TestBearerHeader(t *testing.T) {
	got := httputil.BearerHeader("mytoken")
	if got != "Bearer mytoken" {
		t.Errorf("expected 'Bearer mytoken', got %q", got)
	}
}
