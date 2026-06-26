package domain_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

func TestCompany_DisplaySize(t *testing.T) {
	tests := []struct {
		size     string
		expected string
	}{
		{"1", "Solo founder"},
		{"2-10", "2–10 people"},
		{"11-50", "11–50 people"},
		{"51-200", "51–200 people"},
		{"200+", "200+ people"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		c := domain.Company{Size: domain.CompanySize(tt.size)}
		got := c.DisplaySize()
		if got != tt.expected {
			t.Errorf("DisplaySize(%q) = %q, want %q", tt.size, got, tt.expected)
		}
	}
}

func TestCompany_HasWebsite(t *testing.T) {
	c := domain.Company{Website: "https://example.com"}
	if !c.HasWebsite() {
		t.Fatal("expected HasWebsite() true")
	}
	c2 := domain.Company{}
	if c2.HasWebsite() {
		t.Fatal("expected HasWebsite() false for empty website")
	}
}
