package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

const (
	// SSEChannelLeaderboard is the Redis pub/sub channel for leaderboard updates.
	SSEChannelLeaderboard = "vouch:leaderboard:updated"
	// SSEChannelScore is the Redis pub/sub channel for individual score updates.
	SSEChannelScore = "vouch:score:updated"
)

// SSEHandler streams server-sent events to connected clients.
type SSEHandler struct {
	redis *redis.Client
}

// NewSSEHandler constructs an SSEHandler.
func NewSSEHandler(rdb *redis.Client) *SSEHandler {
	return &SSEHandler{redis: rdb}
}

// scoreEvent is published on SSEChannelScore when a builder's score changes.
type scoreEvent struct {
	BuilderID  string  `json:"builder_id"`
	Username   string  `json:"username"`
	TotalScore float64 `json:"total_score"`
	Tier       string  `json:"tier"`
}

// LeaderboardStream streams leaderboard-updated events to the client via SSE.
// GET /api/v1/sse/leaderboard
func (h *SSEHandler) LeaderboardStream(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	ctx, cancel := context.WithCancel(c.UserContext())
	defer cancel()

	sub := h.redis.Subscribe(ctx, SSEChannelLeaderboard)
	defer sub.Close()

	c.Status(fiber.StatusOK)
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// Send a keep-alive comment every 25 seconds to prevent proxy timeouts.
		ticker := time.NewTicker(25 * time.Second)
		defer ticker.Stop()

		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				fmt.Fprintf(w, "event: leaderboard\ndata: %s\n\n", msg.Payload)
				w.Flush() //nolint:errcheck
			case <-ticker.C:
				fmt.Fprintf(w, ": keep-alive\n\n")
				w.Flush() //nolint:errcheck
			}
		}
	})
	return nil
}

// ScoreStream streams score-updated events for a specific builder via SSE.
// GET /api/v1/sse/scores/:username
func (h *SSEHandler) ScoreStream(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return fiber.ErrBadRequest
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	ctx, cancel := context.WithCancel(c.UserContext())
	defer cancel()

	channel := fmt.Sprintf("%s:%s", SSEChannelScore, username)
	sub := h.redis.Subscribe(ctx, channel)
	defer sub.Close()

	c.Status(fiber.StatusOK)
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		ticker := time.NewTicker(25 * time.Second)
		defer ticker.Stop()

		ch := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				fmt.Fprintf(w, "event: score\ndata: %s\n\n", msg.Payload)
				w.Flush() //nolint:errcheck
			case <-ticker.C:
				fmt.Fprintf(w, ": keep-alive\n\n")
				w.Flush() //nolint:errcheck
			}
		}
	})
	return nil
}

// PublishScoreUpdate publishes a score update event to Redis pub/sub so all
// connected SSE clients are notified in real time. Call after Recalculate.
func PublishScoreUpdate(ctx context.Context, rdb *redis.Client, builderID, username, tier string, total float64) {
	ev := scoreEvent{BuilderID: builderID, Username: username, TotalScore: total, Tier: tier}
	b, _ := json.Marshal(ev)
	// Publish to the per-user channel and the global leaderboard channel.
	rdb.Publish(ctx, fmt.Sprintf("%s:%s", SSEChannelScore, username), b) //nolint:errcheck
	rdb.Publish(ctx, SSEChannelLeaderboard, b)                           //nolint:errcheck
}

