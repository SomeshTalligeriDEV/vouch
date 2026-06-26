// Package observability wires error reporting (Sentry) for the API and worker.
package observability

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

// InitSentry initializes the Sentry client. A blank DSN disables reporting and
// returns a no-op flush, so local development needs no configuration.
func InitSentry(dsn, env, release string) (flush func(), err error) {
	if dsn == "" {
		return func() {}, nil
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      env,
		Release:          release,
		EnableTracing:    false,
		AttachStacktrace: true,
	}); err != nil {
		return func() {}, fmt.Errorf("observability.InitSentry: %w", err)
	}
	return func() { sentry.Flush(2 * time.Second) }, nil
}

// Capture reports an error to Sentry with optional tags. Safe to call even when
// Sentry is disabled (it simply does nothing).
func Capture(err error, tags map[string]string) {
	if err == nil {
		return
	}
	sentry.WithScope(func(scope *sentry.Scope) {
		for k, v := range tags {
			scope.SetTag(k, v)
		}
		sentry.CaptureException(err)
	})
}
