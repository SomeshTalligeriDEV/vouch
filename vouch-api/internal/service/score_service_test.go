package service

import (
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

func TestRecalculate_PersistsScore(t *testing.T) {
	scores := &fakeScoreRepo{}
	projects := newFakeProjectRepo()
	users := newFakeUserRepo()
	stripe := &fakeStripeRepo{}
	enq := &recordingEnqueuer{}

	users.add(&domain.User{ID: "u1", Username: "alice"})
	p := &domain.Project{BuilderID: "u1", Status: domain.ProjectStatusLive, VerifiedUsers: 50}
	projects.Create(ctx(), p)

	svc := NewScoreService(scores, projects, users, stripe, enq)
	score, err := svc.Recalculate(ctx(), "u1")
	if err != nil {
		t.Fatalf("Recalculate: %v", err)
	}
	if score.BuilderID != "u1" {
		t.Fatalf("expected builderID 'u1', got %q", score.BuilderID)
	}
	if score.TotalScore == 0 {
		t.Fatal("expected non-zero score from 50 verified users")
	}
}

func TestRecalculate_EnqueuesScoreUpdate(t *testing.T) {
	enq := &recordingEnqueuer{}
	users := newFakeUserRepo()
	users.add(&domain.User{ID: "u1", Username: "alice"})

	svc := NewScoreService(&fakeScoreRepo{}, newFakeProjectRepo(), users, &fakeStripeRepo{}, enq)
	if err := svc.EnqueueRecalc(ctx(), "u1"); err != nil {
		t.Fatalf("EnqueueRecalc: %v", err)
	}
	if len(enq.scoreCalls) != 1 || enq.scoreCalls[0] != "u1" {
		t.Fatalf("expected one score recalc enqueue for u1, got %v", enq.scoreCalls)
	}
}

func TestRecalculate_StripeMultiplierApplied(t *testing.T) {
	stripe := &fakeStripeRepo{snap: &domain.StripeSnapshot{
		BuilderID:      "u1",
		StripeVerified: true,
		MRR:            200,
	}}
	users := newFakeUserRepo()
	users.add(&domain.User{ID: "u1", Username: "alice"})
	projects := newFakeProjectRepo()
	p := &domain.Project{BuilderID: "u1", Status: domain.ProjectStatusLive, VerifiedUsers: 10}
	projects.Create(ctx(), p)

	svc := NewScoreService(&fakeScoreRepo{}, projects, users, stripe, &recordingEnqueuer{})
	score, err := svc.Recalculate(ctx(), "u1")
	if err != nil {
		t.Fatalf("Recalculate: %v", err)
	}
	if score.StripeMultiplier != 1.0 {
		t.Fatalf("expected Stripe multiplier 1.0 for verified account, got %v", score.StripeMultiplier)
	}
}
