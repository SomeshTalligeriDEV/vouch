// Command seed populates a local MongoDB with demo builders, projects,
// problems, reviews, and computed scores so the app has data to render.
//
// Usage: go run ./cmd/seed   (requires MONGO_URI / MONGO_DB env vars)
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/config"
	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
	repo "github.com/SomeshTalligeriDEV/vouch-api/internal/repository/mongo"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("config load failed")
	}

	ctx := context.Background()
	client, err := repo.Connect(ctx, cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatal().Err(err).Msg("mongo connect failed")
	}
	defer client.Disconnect(ctx)
	if err := client.EnsureIndexes(ctx); err != nil {
		log.Fatal().Err(err).Msg("ensure indexes failed")
	}

	users := repo.NewUserRepo(client)
	projects := repo.NewProjectRepo(client)
	problems := repo.NewProblemRepo(client)
	reviews := repo.NewReviewRepo(client)
	scores := repo.NewScoreRepo(client)
	stripe := repo.NewStripeRepo(client)

	builders := []struct {
		username string
		name     string
		mrr      float64
		users    int
		stripe   bool
		project  string
		tagline  string
	}{
		{"ada", "Ada Lovelace", 4200, 1800, true, "QueryForge", "Visual SQL builder for analysts"},
		{"linus", "Linus T", 0, 320, false, "PatchPilot", "AI PR review for OSS maintainers"},
		{"grace", "Grace H", 12500, 6400, true, "ShipLog", "Changelogs your users actually read"},
	}

	for i, b := range builders {
		u := &domain.User{
			Email:       fmt.Sprintf("%s@example.com", b.username),
			Username:    b.username,
			Name:        b.name,
			Bio:         "Builder shipping real products on Vouch.",
			GitHubID:    int64(1000 + i),
			GitHubLogin: b.username,
			Role:        domain.RoleBuilder,
			IsVerified:  true,
		}
		if err := users.Create(ctx, u); err != nil {
			log.Warn().Err(err).Str("user", b.username).Msg("user exists or failed; skipping")
			existing, gerr := users.GetByUsername(ctx, b.username)
			if gerr != nil {
				continue
			}
			u = existing
		}

		p := &domain.Project{
			BuilderID:     u.ID,
			Title:         b.project,
			Slug:          b.username + "-" + b.project,
			Tagline:       b.tagline,
			Description:   "A real product with paying users, listed on Vouch.",
			LiveURL:       "https://example.com/" + b.username,
			RepoURL:       "https://github.com/" + b.username + "/" + b.project,
			PaymentLink:   "https://buy.stripe.com/demo",
			Tags:          []string{"saas", "developer-tools"},
			Status:        domain.ProjectStatusLive,
			VerifiedUsers: b.users,
			MRR:           b.mrr,
		}
		if err := projects.Create(ctx, p); err != nil {
			log.Warn().Err(err).Msg("project create failed")
		}

		// Two reviews per project from the other builders.
		ratingSum, count := 0, 0
		for j, rb := range builders {
			if rb.username == b.username {
				continue
			}
			rating := 4 + (j % 2)
			rv := &domain.Review{
				ProjectID:        p.ID,
				ReviewerID:       fmt.Sprintf("seed-reviewer-%d-%d", i, j),
				ReviewerUsername: rb.username,
				Rating:           rating,
				Body:             "Genuinely useful. Replaced a tool we were paying for.",
				VerifiedPurchase: true,
			}
			if err := reviews.Create(ctx, rv); err != nil {
				log.Warn().Err(err).Msg("review create failed")
				continue
			}
			ratingSum += rating
			count++
		}
		if count > 0 {
			_ = projects.UpdateRatingStats(ctx, p.ID, domain.ReviewStats{
				Count: count, Average: float64(ratingSum) / float64(count),
			})
			p.ReviewCount = count
			p.AverageRating = float64(ratingSum) / float64(count)
		}

		if b.stripe {
			_ = stripe.Save(ctx, &domain.StripeSnapshot{
				BuilderID: u.ID, MRR: b.mrr, TotalCustomers: b.users, Currency: "usd",
				VerifiedAt: time.Now().UTC(),
			})
		}

		// Compute and persist a score from the seeded inputs.
		score := domain.ComputeScore(domain.ScoreInputs{
			VerifiedUsers:   b.users,
			MRR:             b.mrr,
			AverageRating:   p.AverageRating,
			ReviewCount:     p.ReviewCount,
			NinetyDayGrowth: float64(b.users),
			StripeVerified:  b.stripe,
		})
		score.BuilderID = u.ID
		score.CalculatedAt = time.Now().UTC()
		if err := scores.Upsert(ctx, &score); err != nil {
			log.Warn().Err(err).Msg("score upsert failed")
		}
		log.Info().Str("builder", b.username).Float64("score", score.TotalScore).
			Str("tier", string(score.Tier)).Msg("seeded builder")
	}

	// A couple of open problems on the demand board.
	poster, _ := users.GetByUsername(ctx, "ada")
	if poster != nil {
		for _, pr := range []struct {
			title, desc       string
			budgetMin, budgetMax float64
		}{
			{"Self-hosted status page that ingests Prometheus", "Want a clean status page wired to my existing Prometheus alerts.", 500, 2000},
			{"Stripe → Notion revenue dashboard", "Sync MRR and churn into a Notion database nightly.", 300, 1200},
		} {
			problem := &domain.Problem{
				PosterID:    poster.ID,
				Title:       pr.title,
				Slug:        slugSeed(pr.title),
				Description: pr.desc,
				Tags:        []string{"infra", "saas"},
				BudgetMin:   pr.budgetMin,
				BudgetMax:   pr.budgetMax,
				Status:      domain.ProblemStatusOpen,
				Upvotes:     0,
				UpvotedBy:   []string{},
			}
			if err := problems.Create(ctx, problem); err != nil {
				log.Warn().Err(err).Msg("problem create failed")
			}
		}
	}

	log.Info().Msg("seed complete")
}

// slugSeed is a tiny deterministic slug for seed data.
func slugSeed(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			out = append(out, r)
		case r >= 'A' && r <= 'Z':
			out = append(out, r+32)
		case r == ' ':
			out = append(out, '-')
		}
	}
	return string(out)
}
