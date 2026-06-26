package domain_test

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

func TestProject_IsOwnedBy(t *testing.T) {
	p := domain.Project{BuilderID: "user1"}
	if !p.IsOwnedBy("user1") {
		t.Fatal("expected IsOwnedBy('user1') true")
	}
	if p.IsOwnedBy("user2") {
		t.Fatal("expected IsOwnedBy('user2') false")
	}
}

func TestProject_IsPublic_TrueForLiveStatus(t *testing.T) {
	p := domain.Project{Status: domain.ProjectStatusLive}
	if !p.IsPublic() {
		t.Fatalf("expected IsPublic() true for status %q", p.Status)
	}
}

func TestProject_IsPublic_FalseForDraft(t *testing.T) {
	p := domain.Project{Status: domain.ProjectStatusDraft}
	if p.IsPublic() {
		t.Fatal("expected IsPublic() false for draft")
	}
}
