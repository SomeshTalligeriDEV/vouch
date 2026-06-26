package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// Task type names.
const (
	TypeScoreRecalc         = "score:recalc"
	TypeStripeSync          = "stripe:sync"
	TypeEmailProblemClaimed = "email:problem_claimed"
)

// builderPayload is shared by the score/stripe task types.
type builderPayload struct {
	BuilderID string `json:"builder_id"`
}

// problemPayload identifies a problem for email tasks.
type problemPayload struct {
	ProblemID string `json:"problem_id"`
}

// Enqueuer wraps an Asynq client and implements the service enqueuer ports.
type Enqueuer struct {
	client *asynq.Client
}

// NewEnqueuer constructs an Enqueuer from a Redis connection option.
func NewEnqueuer(redisOpt asynq.RedisConnOpt) *Enqueuer {
	return &Enqueuer{client: asynq.NewClient(redisOpt)}
}

// Close releases the underlying Asynq client.
func (e *Enqueuer) Close() error { return e.client.Close() }

// EnqueueScoreRecalc schedules an async score recalculation.
func (e *Enqueuer) EnqueueScoreRecalc(ctx context.Context, builderID string) error {
	payload, err := json.Marshal(builderPayload{BuilderID: builderID})
	if err != nil {
		return fmt.Errorf("worker.EnqueueScoreRecalc: %w", err)
	}
	task := asynq.NewTask(TypeScoreRecalc, payload, asynq.Queue("default"), asynq.MaxRetry(5))
	if _, err := e.client.EnqueueContext(ctx, task); err != nil {
		return fmt.Errorf("worker.EnqueueScoreRecalc: %w", err)
	}
	return nil
}

// EnqueueStripeSync schedules an async Stripe revenue sync.
func (e *Enqueuer) EnqueueStripeSync(ctx context.Context, builderID string) error {
	payload, err := json.Marshal(builderPayload{BuilderID: builderID})
	if err != nil {
		return fmt.Errorf("worker.EnqueueStripeSync: %w", err)
	}
	task := asynq.NewTask(TypeStripeSync, payload, asynq.Queue("default"), asynq.MaxRetry(5))
	if _, err := e.client.EnqueueContext(ctx, task); err != nil {
		return fmt.Errorf("worker.EnqueueStripeSync: %w", err)
	}
	return nil
}

// EnqueueProblemClaimedEmail schedules the "problem claimed" notification.
func (e *Enqueuer) EnqueueProblemClaimedEmail(ctx context.Context, problemID string) error {
	payload, err := json.Marshal(problemPayload{ProblemID: problemID})
	if err != nil {
		return fmt.Errorf("worker.EnqueueProblemClaimedEmail: %w", err)
	}
	task := asynq.NewTask(TypeEmailProblemClaimed, payload, asynq.Queue("default"), asynq.MaxRetry(3))
	if _, err := e.client.EnqueueContext(ctx, task); err != nil {
		return fmt.Errorf("worker.EnqueueProblemClaimedEmail: %w", err)
	}
	return nil
}
