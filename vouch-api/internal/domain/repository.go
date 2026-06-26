package domain

import "context"

// UserRepository abstracts persistence for users.
type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByGitHubID(ctx context.Context, githubID int64) (*User, error)
	Update(ctx context.Context, u *User) error
	SetStripeAccount(ctx context.Context, id, stripeAccountID string) error
}

// ProjectRepository abstracts persistence for projects.
type ProjectRepository interface {
	Create(ctx context.Context, p *Project) error
	GetByID(ctx context.Context, id string) (*Project, error)
	GetBySlug(ctx context.Context, slug string) (*Project, error)
	List(ctx context.Context, f ProjectFilter) ([]*Project, int64, error)
	ListByBuilder(ctx context.Context, builderID string) ([]*Project, error)
	Update(ctx context.Context, p *Project) error
	UpdateRatingStats(ctx context.Context, projectID string, stats ReviewStats) error
}

// ScoreRepository abstracts persistence for builder scores.
type ScoreRepository interface {
	GetByBuilderID(ctx context.Context, builderID string) (*BuilderScore, error)
	Upsert(ctx context.Context, s *BuilderScore) error
	TopBuilders(ctx context.Context, limit int) ([]*BuilderScore, error)
}

// ProblemRepository abstracts persistence for demand-board problems.
type ProblemRepository interface {
	Create(ctx context.Context, p *Problem) error
	GetByID(ctx context.Context, id string) (*Problem, error)
	GetBySlug(ctx context.Context, slug string) (*Problem, error)
	List(ctx context.Context, f ProblemFilter) ([]*Problem, int64, error)
	Update(ctx context.Context, p *Problem) error
	// Claim atomically transitions an open problem to claimed by builderID.
	Claim(ctx context.Context, id, builderID string) (*Problem, error)
	// AddUpvote atomically records an upvote if the user has not already voted.
	AddUpvote(ctx context.Context, id, userID string) (*Problem, error)
}

// ReviewRepository abstracts persistence for reviews.
type ReviewRepository interface {
	Create(ctx context.Context, r *Review) error
	GetByProjectAndReviewer(ctx context.Context, projectID, reviewerID string) (*Review, error)
	ListByProject(ctx context.Context, projectID string, page, limit int) ([]*Review, int64, error)
	StatsForProject(ctx context.Context, projectID string) (ReviewStats, error)
}

// StripeRepository abstracts persistence for Stripe snapshots.
type StripeRepository interface {
	Save(ctx context.Context, s *StripeSnapshot) error
	Latest(ctx context.Context, builderID string) (*StripeSnapshot, error)
}
