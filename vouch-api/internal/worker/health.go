package worker

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// RunRedisHeartbeat pings Redis every interval and logs failures.
// It exits when ctx is cancelled.
func RunRedisHeartbeat(ctx context.Context, rdb *redis.Client, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := rdb.Ping(ctx).Err(); err != nil {
				log.Error().Err(err).Msg("redis heartbeat failed")
			}
		}
	}
}
