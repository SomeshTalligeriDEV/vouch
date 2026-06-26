package service

import (
	"context"
	"errors"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

func ctx() context.Context { return context.Background() }

// --- Problem claim ---

func TestProblemClaim_SelfClaimRejected(t *testing.T) {
	problems := newFakeProblemRepo()
	enq := &recordingEnqueuer{}
	svc := NewProblemService(problems, enq)

	p, err := svc.Create(ctx(), "user1", ProblemInput{Title: "Need a tool", Description: "x"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, err := svc.Claim(ctx(), p.ID, "user1"); !errors.Is(err, domain.ErrSelfClaim) {
		t.Fatalf("want ErrSelfClaim, got %v", err)
	}
	if len(enq.emailCalls) != 0 {
		t.Fatalf("self-claim should not enqueue email")
	}
}

func TestProblemClaim_SucceedsAndEnqueuesEmail(t *testing.T) {
	problems := newFakeProblemRepo()
	enq := &recordingEnqueuer{}
	svc := NewProblemService(problems, enq)

	p, _ := svc.Create(ctx(), "poster", ProblemInput{Title: "Need a tool", Description: "x"})

	claimed, err := svc.Claim(ctx(), p.ID, "builder")
	if err != nil {
		t.Fatalf("claim: %v", err)
	}
	if claimed.Status != domain.ProblemStatusClaimed || claimed.ClaimedBy != "builder" {
		t.Fatalf("unexpected claimed state: %+v", claimed)
	}
	if len(enq.emailCalls) != 1 || enq.emailCalls[0] != p.ID {
		t.Fatalf("expected one email enqueue for %s, got %v", p.ID, enq.emailCalls)
	}
}

func TestProblemClaim_AlreadyClaimed(t *testing.T) {
	problems := newFakeProblemRepo()
	svc := NewProblemService(problems, &recordingEnqueuer{})
	p, _ := svc.Create(ctx(), "poster", ProblemInput{Title: "Need a tool", Description: "x"})

	if _, err := svc.Claim(ctx(), p.ID, "builderA"); err != nil {
		t.Fatalf("first claim: %v", err)
	}
	if _, err := svc.Claim(ctx(), p.ID, "builderB"); !errors.Is(err, domain.ErrProblemClaimed) {
		t.Fatalf("want ErrProblemClaimed, got %v", err)
	}
}

func TestProblemUpvote_Idempotent(t *testing.T) {
	problems := newFakeProblemRepo()
	svc := NewProblemService(problems, &recordingEnqueuer{})
	p, _ := svc.Create(ctx(), "poster", ProblemInput{Title: "Need a tool", Description: "x"})

	for i := 0; i < 3; i++ {
		if _, err := svc.Upvote(ctx(), p.ID, "voter"); err != nil {
			t.Fatalf("upvote: %v", err)
		}
	}
	got, _ := svc.Get(ctx(), p.ID)
	if got.Upvotes != 1 {
		t.Fatalf("want 1 upvote after repeated votes, got %d", got.Upvotes)
	}
}

// --- Reviews ---

func newReviewSetup(t *testing.T) (*ReviewService, *fakeProjectRepo, *fakeUserRepo, *recordingEnqueuer, *domain.Project) {
	t.Helper()
	projects := newFakeProjectRepo()
	users := newFakeUserRepo()
	reviews := newFakeReviewRepo()
	enq := &recordingEnqueuer{}

	users.add(&domain.User{ID: "builder", Username: "builder"})
	users.add(&domain.User{ID: "reviewer", Username: "reviewer"})
	proj := &domain.Project{BuilderID: "builder", Title: "Tool"}
	_ = projects.Create(ctx(), proj)

	return NewReviewService(reviews, projects, users, enq), projects, users, enq, proj
}

func TestReview_OwnerCannotReview(t *testing.T) {
	svc, _, _, _, proj := newReviewSetup(t)
	_, err := svc.Create(ctx(), "builder", ReviewInput{ProjectID: proj.ID, Rating: 5})
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("want ErrForbidden, got %v", err)
	}
}

func TestReview_InvalidRating(t *testing.T) {
	svc, _, _, _, proj := newReviewSetup(t)
	if _, err := svc.Create(ctx(), "reviewer", ReviewInput{ProjectID: proj.ID, Rating: 0}); !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput for rating 0, got %v", err)
	}
	if _, err := svc.Create(ctx(), "reviewer", ReviewInput{ProjectID: proj.ID, Rating: 6}); !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput for rating 6, got %v", err)
	}
}

func TestReview_DuplicateRejected(t *testing.T) {
	svc, _, _, _, proj := newReviewSetup(t)
	if _, err := svc.Create(ctx(), "reviewer", ReviewInput{ProjectID: proj.ID, Rating: 4}); err != nil {
		t.Fatalf("first review: %v", err)
	}
	if _, err := svc.Create(ctx(), "reviewer", ReviewInput{ProjectID: proj.ID, Rating: 5}); !errors.Is(err, domain.ErrAlreadyExists) {
		t.Fatalf("want ErrAlreadyExists, got %v", err)
	}
}

func TestReview_UpdatesStatsAndEnqueuesRecalc(t *testing.T) {
	svc, projects, _, enq, proj := newReviewSetup(t)
	if _, err := svc.Create(ctx(), "reviewer", ReviewInput{ProjectID: proj.ID, Rating: 4, Body: "good"}); err != nil {
		t.Fatalf("review: %v", err)
	}
	updated := projects.byID[proj.ID]
	if updated.ReviewCount != 1 || updated.AverageRating != 4 {
		t.Fatalf("stats not updated: count=%d avg=%v", updated.ReviewCount, updated.AverageRating)
	}
	if len(enq.scoreCalls) != 1 || enq.scoreCalls[0] != "builder" {
		t.Fatalf("expected recalc enqueue for builder, got %v", enq.scoreCalls)
	}
}

// --- Score recalculation aggregation ---

func TestScoreRecalculate_AggregatesAndAppliesStripe(t *testing.T) {
	users := newFakeUserRepo()
	projects := newFakeProjectRepo()
	scores := &fakeScoreRepo{}
	stripe := &fakeStripeRepo{}
	users.add(&domain.User{ID: "builder", Username: "grace"})

	_ = projects.Create(ctx(), &domain.Project{
		BuilderID: "builder", VerifiedUsers: 100, ReviewCount: 10, AverageRating: 4,
	})
	_ = projects.Create(ctx(), &domain.Project{
		BuilderID: "builder", Status: domain.ProjectStatusArchived, VerifiedUsers: 9999,
	})
	stripe.Save(ctx(), &domain.StripeSnapshot{BuilderID: "builder", MRR: 500, TotalCustomers: 100})

	svc := NewScoreService(scores, projects, users, stripe, &recordingEnqueuer{})
	score, err := svc.Recalculate(ctx(), "builder")
	if err != nil {
		t.Fatalf("recalculate: %v", err)
	}

	// Archived project excluded: users=100 → 1000; revenue 500*2=1000;
	// impact 4*10*5=200; velocity 100*0.1=10; stripe verified ×1.0.
	if score.Breakdown.User != 1000 {
		t.Fatalf("user score = %v, want 1000 (archived excluded)", score.Breakdown.User)
	}
	if score.Breakdown.Revenue != 1000 || !score.StripeVerified || score.StripeMultiplier != 1.0 {
		t.Fatalf("revenue/stripe wrong: %+v", score)
	}
	if scores.upserted == nil {
		t.Fatalf("score was not persisted")
	}
}

func TestScoreRecalculate_UnverifiedMultiplier(t *testing.T) {
	users := newFakeUserRepo()
	projects := newFakeProjectRepo()
	users.add(&domain.User{ID: "b", Username: "linus"})
	_ = projects.Create(ctx(), &domain.Project{BuilderID: "b", VerifiedUsers: 100})

	svc := NewScoreService(&fakeScoreRepo{}, projects, users, &fakeStripeRepo{}, &recordingEnqueuer{})
	score, err := svc.Recalculate(ctx(), "b")
	if err != nil {
		t.Fatalf("recalculate: %v", err)
	}
	// No stripe snapshot → multiplier 0.6; user 100*10=1000 → total 600.
	if score.StripeVerified || score.StripeMultiplier != 0.6 {
		t.Fatalf("expected unverified 0.6 multiplier, got %+v", score)
	}
	if score.TotalScore != 600 {
		t.Fatalf("total = %v, want 600", score.TotalScore)
	}
}
