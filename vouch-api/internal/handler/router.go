package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/middleware"
	"github.com/SomeshTalligeriDEV/vouch-api/pkg/jwt"
)

// Handlers bundles every HTTP handler for registration.
type Handlers struct {
	User    *UserHandler
	Project *ProjectHandler
	Score   *ScoreHandler
	Problem *ProblemHandler
	Review  *ReviewHandler
	Upload  *UploadHandler
}

// Deps carries shared dependencies needed to register routes.
type Deps struct {
	JWT   *jwt.Manager
	Redis *redis.Client
	Log   zerolog.Logger
}

// Register mounts all routes onto the Fiber app under /api/v1.
func Register(app *fiber.App, h Handlers, d Deps) {
	app.Use(middleware.Logger(d.Log))

	auth := middleware.Auth(d.JWT)
	mutationLimiter := middleware.NewRateLimiter(d.Redis, 60, time.Minute).Limit()
	authLimiter := middleware.NewRateLimiter(d.Redis, 20, time.Minute).Limit()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	v1 := app.Group("/api/v1")

	// Auth
	authGrp := v1.Group("/auth")
	authGrp.Post("/github", authLimiter, h.User.GitHubCallback)
	authGrp.Post("/refresh", authLimiter, h.User.Refresh)

	// Users
	users := v1.Group("/users")
	users.Get("/:username", h.User.GetByUsername)
	users.Patch("/me", auth, mutationLimiter, h.User.UpdateMe)
	users.Post("/me/stripe", auth, mutationLimiter, h.User.ConnectStripe)

	// Projects
	projects := v1.Group("/projects")
	projects.Get("/", h.Project.List)
	projects.Post("/", auth, mutationLimiter, h.Project.Create)
	projects.Get("/:id", h.Project.Get)
	projects.Patch("/:id", auth, mutationLimiter, h.Project.Update)
	projects.Delete("/:id", auth, mutationLimiter, h.Project.Delete)

	// Scores
	scores := v1.Group("/scores")
	scores.Get("/", h.Score.Leaderboard)
	scores.Post("/recalculate", auth, mutationLimiter, h.Score.Recalculate)
	scores.Get("/:username", h.Score.GetByUsername)

	// Problems
	problems := v1.Group("/problems")
	problems.Get("/", h.Problem.List)
	problems.Post("/", auth, mutationLimiter, h.Problem.Create)
	problems.Get("/:id", h.Problem.Get)
	problems.Post("/:id/claim", auth, mutationLimiter, h.Problem.Claim)
	problems.Post("/:id/upvote", auth, mutationLimiter, h.Problem.Upvote)

	// Reviews
	reviews := v1.Group("/reviews")
	reviews.Post("/", auth, mutationLimiter, h.Review.Create)
	reviews.Get("/project/:id", h.Review.ListByProject)

	// Uploads
	uploads := v1.Group("/uploads")
	uploads.Post("/presign", auth, mutationLimiter, h.Upload.Presign)
}
