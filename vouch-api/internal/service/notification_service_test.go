package service

import (
	"context"
	"testing"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

type recordingNotifier struct {
	sent []struct{ to, subject, body string }
}

func (n *recordingNotifier) Send(_ context.Context, to, subject, body string) error {
	n.sent = append(n.sent, struct{ to, subject, body string }{to, subject, body})
	return nil
}

func TestNotifyProblemClaimed_SendsEmailToPoster(t *testing.T) {
	users := newFakeUserRepo()
	poster := users.add(&domain.User{ID: "poster1", Username: "acme", Email: "acme@example.com"})
	builder := users.add(&domain.User{ID: "builder1", Username: "alice"})

	problems := newFakeProblemRepo()
	prob := &domain.Problem{
		PosterID:  poster.ID,
		ClaimedBy: builder.ID,
		Title:     "Need a Slack bot",
		Status:    domain.ProblemStatusClaimed,
	}
	_ = problems.Create(ctx(), prob)

	notifier := &recordingNotifier{}
	svc := NewNotificationService(users, problems, notifier, "http://localhost:3000")

	if err := svc.NotifyProblemClaimed(context.Background(), prob.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(notifier.sent) != 1 {
		t.Fatalf("expected 1 email sent, got %d", len(notifier.sent))
	}
	if notifier.sent[0].to != "acme@example.com" {
		t.Errorf("expected email to acme@example.com, got %s", notifier.sent[0].to)
	}
	if notifier.sent[0].subject == "" {
		t.Error("expected non-empty email subject")
	}
}

func TestNotifyProblemClaimed_NoPosterEmail_NoSend(t *testing.T) {
	users := newFakeUserRepo()
	poster := users.add(&domain.User{ID: "poster2", Username: "anon", Email: ""}) // no email
	_ = poster

	problems := newFakeProblemRepo()
	prob := &domain.Problem{
		PosterID: "poster2",
		Title:    "Problem with no email poster",
		Status:   domain.ProblemStatusClaimed,
	}
	_ = problems.Create(ctx(), prob)

	notifier := &recordingNotifier{}
	svc := NewNotificationService(users, problems, notifier, "http://localhost:3000")

	if err := svc.NotifyProblemClaimed(context.Background(), prob.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifier.sent) != 0 {
		t.Errorf("expected no emails sent for poster with empty email, got %d", len(notifier.sent))
	}
}
