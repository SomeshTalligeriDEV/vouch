package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/config"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/external"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/handler"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/observability"
	repo "github.com/SomeshTalligeriDEV/vouch-api/internal/repository/mongo"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/service"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/worker"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/jwt"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/response"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/validator"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("config load failed")
	}

	logger := newLogger(cfg)

	flushSentry, err := observability.InitSentry(cfg.SentryDSN, cfg.Env, "vouch-api")
	if err != nil {
		logger.Fatal().Err(err).Msg("sentry init failed")
	}
	defer flushSentry()

	ctx := context.Background()
	mongoClient, err := repo.Connect(ctx, cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		logger.Fatal().Err(err).Msg("mongo connect failed")
	}
	defer mongoClient.Disconnect(ctx)
	if err := mongoClient.EnsureIndexes(ctx); err != nil {
		logger.Fatal().Err(err).Msg("ensure indexes failed")
	}

	redisOpt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("redis url parse failed")
	}
	rdb := redis.NewClient(redisOpt)
	defer rdb.Close()

	asynqOpt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("asynq redis uri parse failed")
	}
	enqueuer := worker.NewEnqueuer(asynqOpt)
	defer enqueuer.Close()

	// Repositories
	userRepo := repo.NewUserRepo(mongoClient)
	projectRepo := repo.NewProjectRepo(mongoClient)
	scoreRepo := repo.NewScoreRepo(mongoClient)
	problemRepo := repo.NewProblemRepo(mongoClient)
	reviewRepo := repo.NewReviewRepo(mongoClient)
	stripeRepo := repo.NewStripeRepo(mongoClient)
	companyRepo := repo.NewCompanyRepo(mongoClient)

	// External gateways
	githubClient := external.NewGitHubClient(cfg.GitHubClientID, cfg.GitHubClientSecret, cfg.GitHubRedirectURL)
	stripeClient := external.NewStripeClient(cfg.StripeClientID, cfg.StripeSecretKey)
	r2Presigner := external.NewR2Presigner(cfg.R2AccessKey, cfg.R2SecretKey, cfg.R2Endpoint, cfg.R2Bucket, cfg.R2PublicURL)

	// Core
	jwtMgr := jwt.NewManager(cfg.JWTSecret, cfg.JWTRefreshSecret)
	val := validator.New()

	// Services
	userSvc := service.NewUserService(userRepo, jwtMgr, githubClient)
	projectSvc := service.NewProjectService(projectRepo, enqueuer)
	scoreSvc := service.NewScoreService(scoreRepo, projectRepo, userRepo, stripeRepo, enqueuer)
	problemSvc := service.NewProblemService(problemRepo, enqueuer)
	reviewSvc := service.NewReviewService(reviewRepo, projectRepo, userRepo, enqueuer)
	stripeSvc := service.NewStripeService(userRepo, stripeRepo, stripeClient, enqueuer)
	uploadSvc := service.NewUploadService(r2Presigner)
	companySvc := service.NewCompanyService(companyRepo, jwtMgr)
	adminSvc := service.NewAdminService(userRepo, companyRepo, projectRepo, problemRepo, reviewRepo)

	// HTTP
	app := fiber.New(fiber.Config{
		AppName:               "vouch-api",
		DisableStartupMessage: true,
		BodyLimit:             2 * 1024 * 1024, // 2 MB — presigned uploads bypass this, raw JSON only
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Error().Err(err).Str("path", c.Path()).Msg("unhandled error")
			observability.Capture(err, map[string]string{
				"path":   c.Path(),
				"method": c.Method(),
			})
			return response.Error(c, fiber.StatusInternalServerError, "internal_error", "something went wrong")
		},
	})

	handler.Register(app, handler.Handlers{
		User:    handler.NewUserHandler(userSvc, stripeSvc, val),
		Project: handler.NewProjectHandler(projectSvc, val),
		Score:   handler.NewScoreHandler(scoreSvc),
		Problem: handler.NewProblemHandler(problemSvc, val),
		Review:  handler.NewReviewHandler(reviewSvc, val),
		Upload:  handler.NewUploadHandler(uploadSvc, val),
		Company: handler.NewCompanyHandler(companySvc, val),
		Admin:   handler.NewAdminHandler(adminSvc),
	}, handler.Deps{JWT: jwtMgr, Redis: rdb, Log: logger, AllowedOrigins: cfg.AllowedOrigins})

	// Graceful shutdown
	go func() {
		addr := ":" + cfg.Port
		logger.Info().Str("addr", addr).Str("env", cfg.Env).Msg("vouch-api listening")
		if err := app.Listen(addr); err != nil {
			logger.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info().Msg("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("graceful shutdown failed")
	}
}

func newLogger(cfg *config.Config) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	if cfg.IsProduction() {
		return zerolog.New(os.Stdout).With().Timestamp().Str("service", "vouch-api").Logger()
	}
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen}).
		With().Timestamp().Logger()
}
