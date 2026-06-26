package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	Company *CompanyHandler
	Admin   *AdminHandler
}

// Deps carries shared dependencies needed to register routes.
type Deps struct {
	JWT        *jwt.Manager
	Redis      *redis.Client
	Log        zerolog.Logger
	AllowedOrigins string // comma-separated origins, e.g. "https://vouch.dev,https://www.vouch.dev"
}

// Register mounts all routes onto the Fiber app under /api/v1.
func Register(app *fiber.App, h Handlers, d Deps) {
	// CORS must be first so preflight OPTIONS requests are handled before auth.
	app.Use(cors.New(cors.Config{
		AllowOrigins:     d.AllowedOrigins,
		AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Authorization,X-Request-ID",
		ExposeHeaders:    "X-Request-ID",
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(d.Log))

	auth := middleware.Auth(d.JWT)
	mutationLimiter := middleware.NewRateLimiter(d.Redis, 60, time.Minute).Limit()
	authLimiter := middleware.NewRateLimiter(d.Redis, 20, time.Minute).Limit()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	app.Get("/ready", func(c *fiber.Ctx) error {
		if err := d.Redis.Ping(c.UserContext()).Err(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "not_ready", "error": "redis"})
		}
		return c.JSON(fiber.Map{"status": "ready"})
	})

	v1 := app.Group("/api/v1")

	// Auth
	authGrp := v1.Group("/auth")
	authGrp.Post("/github", authLimiter, h.User.GitHubCallback)
	authGrp.Post("/refresh", authLimiter, h.User.Refresh)
	authGrp.Post("/logout", h.User.Logout)

	companyOnly := middleware.RequireSubjectType("company")
	userOnly := middleware.RequireSubjectType("user")

	// Users
	users := v1.Group("/users")
	users.Get("/:username", h.User.GetByUsername)
	users.Get("/me", auth, userOnly, h.User.GetMe)
	users.Patch("/me", auth, userOnly, mutationLimiter, h.User.UpdateMe)
	users.Post("/me/stripe", auth, userOnly, mutationLimiter, h.User.ConnectStripe)

	// Projects
	projects := v1.Group("/projects")
	projects.Get("/", h.Project.List)
	projects.Post("/", auth, userOnly, mutationLimiter, h.Project.Create)
	projects.Get("/:id", h.Project.Get)
	projects.Patch("/:id", auth, userOnly, mutationLimiter, h.Project.Update)
	projects.Delete("/:id", auth, userOnly, mutationLimiter, h.Project.Delete)

	// Scores
	scores := v1.Group("/scores")
	scores.Get("/", h.Score.Leaderboard)
	scores.Post("/recalculate", auth, userOnly, mutationLimiter, h.Score.Recalculate)
	scores.Get("/:username", h.Score.GetByUsername)

	// Problems — both builders and companies can post/claim/upvote
	problems := v1.Group("/problems")
	problems.Get("/", h.Problem.List)
	problems.Post("/", auth, mutationLimiter, h.Problem.Create)
	problems.Get("/:id", h.Problem.Get)
	problems.Post("/:id/claim", auth, userOnly, mutationLimiter, h.Problem.Claim)
	problems.Post("/:id/upvote", auth, mutationLimiter, h.Problem.Upvote)

	// Reviews — builders only (they're reviewing projects they used/paid for)
	reviews := v1.Group("/reviews")
	reviews.Post("/", auth, userOnly, mutationLimiter, h.Review.Create)
	reviews.Get("/project/:id", h.Review.ListByProject)

	// Uploads
	uploads := v1.Group("/uploads")
	uploads.Post("/presign", auth, mutationLimiter, h.Upload.Presign)

	// Companies (email + password auth)
	companies := v1.Group("/companies")
	companies.Post("/register", authLimiter, h.Company.Register)
	companies.Post("/login", authLimiter, h.Company.Login)
	companies.Post("/refresh", authLimiter, h.Company.Refresh)
	companies.Post("/logout", h.Company.Logout)
	companies.Get("/me", auth, companyOnly, h.Company.GetMe)
	companies.Patch("/me", auth, companyOnly, mutationLimiter, h.Company.UpdateMe)
	companies.Get("/:slug", h.Company.GetBySlug)

	// Admin (requires auth + admin role)
	adminGrp := v1.Group("/admin", auth, middleware.RequireRole("admin"))
	adminGrp.Get("/stats", h.Admin.Stats)
	adminGrp.Get("/companies", h.Admin.ListCompanies)

	// SSE — real-time streams (no auth required; score data is public)
	sse := NewSSEHandler(d.Redis)
	sseGrp := v1.Group("/sse")
	sseGrp.Get("/leaderboard", sse.LeaderboardStream)
	sseGrp.Get("/scores/:username", sse.ScoreStream)
}
