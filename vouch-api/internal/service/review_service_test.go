package service

import (
	"errors"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

func TestCreateReview_DuplicateRejected(t *testing.T) {
	projects := newFakeProjectRepo()
	users := newFakeUserRepo()
	enq := &recordingEnqueuer{}

	// Seed a project.
	p := &domain.Project{BuilderID: "builder1", Title: "Cool App", Status: domain.ProjectStatusLive}
	projects.Create(ctx(), p)

	// Seed the reviewer user.
	reviewer := &domain.User{ID: "reviewer1", Username: "dave"}
	users.add(reviewer)

	reviews := newFakeReviewRepo()
	svc := NewReviewService(reviews, projects, users, enq)

	_, err := svc.Create(ctx(), "reviewer1", ReviewInput{
		ProjectID: p.ID,
		Rating:    5,
		Body:      "Excellent product",
	})
	if err != nil {
		t.Fatalf("first review: %v", err)
	}

	// Second review from same user should fail.
	_, err = svc.Create(ctx(), "reviewer1", ReviewInput{
		ProjectID: p.ID,
		Rating:    4,
		Body:      "Changed my mind",
	})
	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists on duplicate review, got %v", err)
	}
}

func TestCreateReview_OwnerCannotReviewOwnProject(t *testing.T) {
	projects := newFakeProjectRepo()
	users := newFakeUserRepo()
	enq := &recordingEnqueuer{}

	p := &domain.Project{BuilderID: "builder1", Title: "My App", Status: domain.ProjectStatusLive}
	projects.Create(ctx(), p)

	builder := &domain.User{ID: "builder1", Username: "builder"}
	users.add(builder)

	reviews := newFakeReviewRepo()
	svc := NewReviewService(reviews, projects, users, enq)

	_, err := svc.Create(ctx(), "builder1", ReviewInput{
		ProjectID: p.ID,
		Rating:    5,
		Body:      "My own project is great",
	})
	if err == nil {
		t.Fatal("expected error when owner reviews own project")
	}
}
