package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/handler"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
)

// ScoreWorker handles async score recalculation tasks.
type ScoreWorker struct {
	scores *service.ScoreService
	rdb    *redis.Client
}

// NewScoreWorker constructs a ScoreWorker.
func NewScoreWorker(scores *service.ScoreService, rdb *redis.Client) *ScoreWorker {
	return &ScoreWorker{scores: scores, rdb: rdb}
}

// Handle processes a single score:recalc task and publishes an SSE event.
func (w *ScoreWorker) Handle(ctx context.Context, t *asynq.Task) error {
	var p builderPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("ScoreWorker.Handle unmarshal: %w: %v", asynq.SkipRetry, err)
	}
	score, err := w.scores.Recalculate(ctx, p.BuilderID)
	if err != nil {
		return fmt.Errorf("ScoreWorker.Handle: %w", err)
	}
	// Publish real-time update so SSE clients see the new score immediately.
	entry, _ := w.scores.LeaderboardEntryForBuilder(ctx, p.BuilderID)
	username := ""
	if entry != nil {
		username = entry.Username
	}
	handler.PublishScoreUpdate(ctx, w.rdb, p.BuilderID, username, string(score.Tier), score.TotalScore)
	return nil
}
