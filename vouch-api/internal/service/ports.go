package service

import "context"

// ScoreEnqueuer lets services schedule an asynchronous score recalculation
// without depending on the queue implementation. The worker package supplies
// a concrete implementation backed by Asynq.
type ScoreEnqueuer interface {
	EnqueueScoreRecalc(ctx context.Context, builderID string) error
}

// StripeSyncEnqueuer lets services schedule an asynchronous Stripe data sync.
type StripeSyncEnqueuer interface {
	EnqueueStripeSync(ctx context.Context, builderID string) error
}

// EmailEnqueuer lets services schedule asynchronous transactional emails.
type EmailEnqueuer interface {
	EnqueueProblemClaimedEmail(ctx context.Context, problemID string) error
}

// Notifier abstracts the transactional email provider (Resend).
type Notifier interface {
	Send(ctx context.Context, to, subject, html string) error
}
