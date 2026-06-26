package main

import (
	"context"
	"time"

	"os"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/config"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/external"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/observability"
	repo "github.com/SomeshTalligeriDEV/vouch-api/internal/repository/mongo"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("config load failed")
	}
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen}).
		With().Timestamp().Str("service", "vouch-worker").Logger()

	ctx := context.Background()
	mongoClient, err := repo.Connect(ctx, cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		logger.Fatal().Err(err).Msg("mongo connect failed")
	}
	defer mongoClient.Disconnect(ctx)

	asynqOpt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("asynq redis uri parse failed")
	}

	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("redis url parse failed")
	}
	rdb := redis.NewClient(redisOpt)
	defer rdb.Close()

	// Self-enqueuer so services can chain follow-up tasks.
	enqueuer := worker.NewEnqueuer(asynqOpt)
	defer enqueuer.Close()

	// Repositories
	userRepo := repo.NewUserRepo(mongoClient)
	projectRepo := repo.NewProjectRepo(mongoClient)
	scoreRepo := repo.NewScoreRepo(mongoClient)
	problemRepo := repo.NewProblemRepo(mongoClient)
	stripeRepo := repo.NewStripeRepo(mongoClient)

	stripeClient := external.NewStripeClient(cfg.StripeClientID, cfg.StripeSecretKey)
	resendClient := external.NewResendClient(cfg.ResendAPIKey, cfg.EmailFrom)

	// Services needed by workers
	scoreSvc := service.NewScoreService(scoreRepo, projectRepo, userRepo, stripeRepo, enqueuer)
	stripeSvc := service.NewStripeService(userRepo, stripeRepo, stripeClient, enqueuer)
	notificationSvc := service.NewNotificationService(userRepo, problemRepo, resendClient, cfg.AppURL)

	scoreWorker := worker.NewScoreWorker(scoreSvc, rdb)
	stripeWorker := worker.NewStripeWorker(stripeSvc)
	emailWorker := worker.NewEmailWorker(notificationSvc)

	flushSentry, err := observability.InitSentry(cfg.SentryDSN, cfg.Env, "vouch-worker")
	if err != nil {
		logger.Fatal().Err(err).Msg("sentry init failed")
	}
	defer flushSentry()

	srv := asynq.NewServer(asynqOpt, asynq.Config{
		Concurrency: 10,
		Queues:      map[string]int{"default": 10},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			logger.Error().Err(err).Str("task", task.Type()).Msg("task failed")
			observability.Capture(err, map[string]string{"task": task.Type()})
		}),
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc(worker.TypeScoreRecalc, scoreWorker.Handle)
	mux.HandleFunc(worker.TypeStripeSync, stripeWorker.Handle)
	mux.HandleFunc(worker.TypeEmailProblemClaimed, emailWorker.HandleProblemClaimed)

	logger.Info().Msg("vouch-worker starting")
	if err := srv.Run(mux); err != nil {
		logger.Fatal().Err(err).Msg("worker server failed")
	}
}
