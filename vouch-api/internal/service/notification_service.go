package service

import (
	"context"
	"fmt"
	"html"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// NotificationService builds and sends transactional emails. It is invoked
// from the worker so a slow email provider never blocks an API request.
type NotificationService struct {
	users    domain.UserRepository
	problems domain.ProblemRepository
	notifier Notifier
	appURL   string
}

// NewNotificationService constructs a NotificationService.
func NewNotificationService(
	users domain.UserRepository,
	problems domain.ProblemRepository,
	notifier Notifier,
	appURL string,
) *NotificationService {
	return &NotificationService{users: users, problems: problems, notifier: notifier, appURL: appURL}
}

// NotifyProblemClaimed emails the poster that a builder claimed their problem.
func (s *NotificationService) NotifyProblemClaimed(ctx context.Context, problemID string) error {
	problem, err := s.problems.GetByID(ctx, problemID)
	if err != nil {
		return fmt.Errorf("NotificationService.NotifyProblemClaimed: %w", err)
	}
	poster, err := s.users.GetByID(ctx, problem.PosterID)
	if err != nil {
		return fmt.Errorf("NotificationService.NotifyProblemClaimed: %w", err)
	}
	if poster.Email == "" {
		return nil // nothing to send to
	}
	builderName := "A builder"
	if problem.ClaimedBy != "" {
		if b, err := s.users.GetByID(ctx, problem.ClaimedBy); err == nil {
			builderName = "@" + b.Username
		}
	}

	subject := fmt.Sprintf("%s claimed your problem on Vouch", builderName)
	link := fmt.Sprintf("%s/problems/%s", s.appURL, problem.ID)
	body := fmt.Sprintf(`
		<div style="font-family:system-ui,sans-serif;max-width:480px">
		  <h2>Your problem just got claimed 🎉</h2>
		  <p><strong>%s</strong> claimed your problem <strong>%s</strong> and is shipping a solution.</p>
		  <p>You're in line to become their first paying user.</p>
		  <p><a href="%s" style="color:#6366f1">View the problem →</a></p>
		  <hr/>
		  <p style="color:#888;font-size:12px">Vouch — verified proof of what you've built.</p>
		</div>`,
		html.EscapeString(builderName), html.EscapeString(problem.Title), link)

	if err := s.notifier.Send(ctx, poster.Email, subject, body); err != nil {
		return fmt.Errorf("NotificationService.NotifyProblemClaimed: %w", err)
	}
	return nil
}
