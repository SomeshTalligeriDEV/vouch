package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
)

// StripeWorker handles async Stripe revenue sync tasks.
type StripeWorker struct {
	stripe *service.StripeService
}

// NewStripeWorker constructs a StripeWorker.
func NewStripeWorker(stripe *service.StripeService) *StripeWorker {
	return &StripeWorker{stripe: stripe}
}

// Handle processes a single stripe:sync task.
func (w *StripeWorker) Handle(ctx context.Context, t *asynq.Task) error {
	var p builderPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("StripeWorker.Handle unmarshal: %w: %v", asynq.SkipRetry, err)
	}
	if _, err := w.stripe.Sync(ctx, p.BuilderID); err != nil {
		return fmt.Errorf("StripeWorker.Handle: %w", err)
	}
	return nil
}
