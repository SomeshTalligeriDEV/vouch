package service

import (
	"errors"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

func TestCreateProblem_Success(t *testing.T) {
	problems := newFakeProblemRepo()
	users := newFakeUserRepo()
	enq := &recordingEnqueuer{}

	poster := &domain.User{ID: "poster1", Username: "alice"}
	users.add(poster)

	svc := NewProblemService(problems, enq)
	p, err := svc.Create(ctx(), "poster1", ProblemInput{
		Title:       "Need a Stripe integration",
		Description: "Looking for a builder to audit our checkout flow.",
		Tags:        []string{"stripe", "payments"},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if p.Title != "Need a Stripe integration" {
		t.Fatalf("unexpected title: %q", p.Title)
	}
	if p.Status != domain.ProblemStatusOpen {
		t.Fatalf("expected open status, got %q", p.Status)
	}
}

func TestCreateProblem_BudgetMaxLessThanMin_RejectsInput(t *testing.T) {
	svc := NewProblemService(newFakeProblemRepo(), &recordingEnqueuer{})
	_, err := svc.Create(ctx(), "poster1", ProblemInput{
		Title:       "Some problem",
		Description: "Description here.",
		BudgetMin:   500,
		BudgetMax:   100,
	})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for BudgetMax < BudgetMin, got %v", err)
	}
}
