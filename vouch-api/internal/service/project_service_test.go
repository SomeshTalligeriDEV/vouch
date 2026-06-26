package service

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

func TestCreateProject_Success(t *testing.T) {
	projects := newFakeProjectRepo()
	users := newFakeUserRepo()
	enq := &recordingEnqueuer{}

	users.add(&domain.User{ID: "u1", Username: "alice"})

	svc := NewProjectService(projects, enq)
	p, err := svc.Create(ctx(), "u1", CreateInput{
		Title:       "QueryForge",
		Description: "Visual SQL builder",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if p.BuilderID != "u1" {
		t.Fatalf("expected BuilderID 'u1', got %q", p.BuilderID)
	}
}

func TestCreateProject_SlugDerived(t *testing.T) {
	svc := NewProjectService(newFakeProjectRepo(), &recordingEnqueuer{})
	p, err := svc.Create(ctx(), "u1", CreateInput{Title: "My Cool SaaS"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if p.Slug == "" {
		t.Fatal("expected non-empty slug")
	}
	if p.BuilderID != "u1" {
		t.Fatalf("expected builderID 'u1', got %q", p.BuilderID)
	}
}
