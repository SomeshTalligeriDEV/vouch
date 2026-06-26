package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// ProjectRepo is the MongoDB implementation of domain.ProjectRepository.
type ProjectRepo struct {
	coll *mongo.Collection
}

// NewProjectRepo constructs a ProjectRepo.
func NewProjectRepo(c *Client) *ProjectRepo {
	return &ProjectRepo{coll: c.DB().Collection("projects")}
}

func (r *ProjectRepo) Create(ctx context.Context, p *domain.Project) error {
	if p.ID == "" {
		p.ID = primitive.NewObjectID().Hex()
	}
	now := time.Now().UTC()
	p.CreatedAt, p.UpdatedAt = now, now
	if _, err := r.coll.InsertOne(ctx, p); err != nil {
		return fmt.Errorf("ProjectRepo.Create: %w", mapMongoErr(err))
	}
	return nil
}

func (r *ProjectRepo) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	var p domain.Project
	if err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&p); err != nil {
		return nil, fmt.Errorf("ProjectRepo.GetByID: %w", mapMongoErr(err))
	}
	return &p, nil
}

func (r *ProjectRepo) GetBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	var p domain.Project
	if err := r.coll.FindOne(ctx, bson.M{"slug": slug}).Decode(&p); err != nil {
		return nil, fmt.Errorf("ProjectRepo.GetBySlug: %w", mapMongoErr(err))
	}
	return &p, nil
}

func (r *ProjectRepo) List(ctx context.Context, f domain.ProjectFilter) ([]*domain.Project, int64, error) {
	filter := bson.M{}
	if f.BuilderID != "" {
		filter["builder_id"] = f.BuilderID
	}
	if f.Status != "" {
		filter["status"] = f.Status
	}
	if f.ForSale != nil {
		filter["for_sale"] = *f.ForSale
	}
	if f.Tag != "" {
		filter["tags"] = f.Tag
	}
	if f.Search != "" {
		filter["title"] = bson.M{"$regex": f.Search, "$options": "i"}
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("ProjectRepo.List count: %w", mapMongoErr(err))
	}

	opts := options.Find().
		SetSort(sortForProjects(f.SortBy)).
		SetSkip(int64((f.Page - 1) * f.Limit)).
		SetLimit(int64(f.Limit))

	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("ProjectRepo.List find: %w", mapMongoErr(err))
	}
	var out []*domain.Project
	if err := cur.All(ctx, &out); err != nil {
		return nil, 0, fmt.Errorf("ProjectRepo.List decode: %w", mapMongoErr(err))
	}
	return out, total, nil
}

func sortForProjects(sortBy string) bson.D {
	switch sortBy {
	case "mrr":
		return bson.D{{Key: "mrr", Value: -1}}
	case "users":
		return bson.D{{Key: "verified_users", Value: -1}}
	case "rating":
		return bson.D{{Key: "average_rating", Value: -1}}
	default:
		return bson.D{{Key: "created_at", Value: -1}}
	}
}

func (r *ProjectRepo) ListByBuilder(ctx context.Context, builderID string) ([]*domain.Project, error) {
	cur, err := r.coll.Find(ctx, bson.M{"builder_id": builderID})
	if err != nil {
		return nil, fmt.Errorf("ProjectRepo.ListByBuilder: %w", mapMongoErr(err))
	}
	var out []*domain.Project
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("ProjectRepo.ListByBuilder decode: %w", mapMongoErr(err))
	}
	return out, nil
}

func (r *ProjectRepo) Update(ctx context.Context, p *domain.Project) error {
	p.UpdatedAt = time.Now().UTC()
	set := bson.M{
		"title":          p.Title,
		"tagline":        p.Tagline,
		"description":    p.Description,
		"logo_url":       p.LogoURL,
		"live_url":       p.LiveURL,
		"repo_url":       p.RepoURL,
		"payment_link":   p.PaymentLink,
		"tags":           p.Tags,
		"status":         p.Status,
		"for_sale":       p.ForSale,
		"ask_price":      p.AskPrice,
		"verified_users": p.VerifiedUsers,
		"mrr":            p.MRR,
		"updated_at":     p.UpdatedAt,
	}
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": p.ID}, bson.M{"$set": set})
	if err != nil {
		return fmt.Errorf("ProjectRepo.Update: %w", mapMongoErr(err))
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("ProjectRepo.Update: %w", domain.ErrNotFound)
	}
	return nil
}

func (r *ProjectRepo) UpdateRatingStats(ctx context.Context, projectID string, stats domain.ReviewStats) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": projectID}, bson.M{"$set": bson.M{
		"review_count":   stats.Count,
		"average_rating": stats.Average,
		"updated_at":     time.Now().UTC(),
	}})
	if err != nil {
		return fmt.Errorf("ProjectRepo.UpdateRatingStats: %w", mapMongoErr(err))
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("ProjectRepo.UpdateRatingStats: %w", domain.ErrNotFound)
	}
	return nil
}
