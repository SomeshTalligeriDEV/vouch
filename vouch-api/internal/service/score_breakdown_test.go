package service

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

func TestScoreBreakdown_ZeroProjectsZeroScore(t *testing.T) {
	users := newFakeUserRepo()
	users.add(&domain.User{ID: "u1", Username: "alice"})

	svc := NewScoreService(&fakeScoreRepo{}, newFakeProjectRepo(), users, &fakeStripeRepo{}, &recordingEnqueuer{})
	score, err := svc.Recalculate(ctx(), "u1")
	if err != nil {
		t.Fatalf("Recalculate: %v", err)
	}
	if score.TotalScore != 0 {
		t.Fatalf("expected 0 score with no projects, got %v", score.TotalScore)
	}
}

func TestScoreBreakdown_MultipleProjects(t *testing.T) {
	users := newFakeUserRepo()
	users.add(&domain.User{ID: "u2", Username: "bob"})
	projects := newFakeProjectRepo()

	for i := 0; i < 3; i++ {
		p := &domain.Project{
			BuilderID:     "u2",
			Status:        domain.ProjectStatusLive,
			VerifiedUsers: 20,
		}
		projects.Create(ctx(), p)
	}

	svc := NewScoreService(&fakeScoreRepo{}, projects, users, &fakeStripeRepo{}, &recordingEnqueuer{})
	score, err := svc.Recalculate(ctx(), "u2")
	if err != nil {
		t.Fatalf("Recalculate: %v", err)
	}
	if score.TotalScore == 0 {
		t.Fatal("expected non-zero score for 3 live projects")
	}
}
