package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
)

// ScoreWorker handles async score recalculation tasks.
type ScoreWorker struct {
	scores *service.ScoreService
}

// NewScoreWorker constructs a ScoreWorker.
func NewScoreWorker(scores *service.ScoreService) *ScoreWorker {
	return &ScoreWorker{scores: scores}
}

// Handle processes a single score:recalc task.
func (w *ScoreWorker) Handle(ctx context.Context, t *asynq.Task) error {
	var p builderPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		// Unrecoverable: a malformed payload will never succeed on retry.
		return fmt.Errorf("ScoreWorker.Handle unmarshal: %w: %v", asynq.SkipRetry, err)
	}
	if _, err := w.scores.Recalculate(ctx, p.BuilderID); err != nil {
		return fmt.Errorf("ScoreWorker.Handle: %w", err)
	}
	return nil
}
