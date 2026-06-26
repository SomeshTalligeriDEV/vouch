package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
)

// EmailWorker handles async transactional email tasks.
type EmailWorker struct {
	notifications *service.NotificationService
}

// NewEmailWorker constructs an EmailWorker.
func NewEmailWorker(notifications *service.NotificationService) *EmailWorker {
	return &EmailWorker{notifications: notifications}
}

// HandleProblemClaimed processes an email:problem_claimed task.
func (w *EmailWorker) HandleProblemClaimed(ctx context.Context, t *asynq.Task) error {
	var p problemPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("EmailWorker.HandleProblemClaimed unmarshal: %w: %v", asynq.SkipRetry, err)
	}
	if err := w.notifications.NotifyProblemClaimed(ctx, p.ProblemID); err != nil {
		return fmt.Errorf("EmailWorker.HandleProblemClaimed: %w", err)
	}
	return nil
}
