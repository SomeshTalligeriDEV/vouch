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

// ProblemRepo is the MongoDB implementation of domain.ProblemRepository.
type ProblemRepo struct {
	coll *mongo.Collection
}

// NewProblemRepo constructs a ProblemRepo.
func NewProblemRepo(c *Client) *ProblemRepo {
	return &ProblemRepo{coll: c.DB().Collection("problems")}
}

func (r *ProblemRepo) Create(ctx context.Context, p *domain.Problem) error {
	if p.ID == "" {
		p.ID = primitive.NewObjectID().Hex()
	}
	if p.UpvotedBy == nil {
		p.UpvotedBy = []string{}
	}
	now := time.Now().UTC()
	p.CreatedAt, p.UpdatedAt = now, now
	if _, err := r.coll.InsertOne(ctx, p); err != nil {
		return fmt.Errorf("ProblemRepo.Create: %w", mapMongoErr(err))
	}
	return nil
}

func (r *ProblemRepo) GetByID(ctx context.Context, id string) (*domain.Problem, error) {
	var p domain.Problem
	if err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&p); err != nil {
		return nil, fmt.Errorf("ProblemRepo.GetByID: %w", mapMongoErr(err))
	}
	return &p, nil
}

func (r *ProblemRepo) GetBySlug(ctx context.Context, slug string) (*domain.Problem, error) {
	var p domain.Problem
	if err := r.coll.FindOne(ctx, bson.M{"slug": slug}).Decode(&p); err != nil {
		return nil, fmt.Errorf("ProblemRepo.GetBySlug: %w", mapMongoErr(err))
	}
	return &p, nil
}

func (r *ProblemRepo) List(ctx context.Context, f domain.ProblemFilter) ([]*domain.Problem, int64, error) {
	filter := bson.M{}
	if f.PosterID != "" {
		filter["poster_id"] = f.PosterID
	}
	if f.ClaimedBy != "" {
		filter["claimed_by"] = f.ClaimedBy
	}
	if f.Status != "" {
		filter["status"] = f.Status
	}
	if f.Tag != "" {
		filter["tags"] = f.Tag
	}
	if f.Search != "" {
		filter["title"] = bson.M{"$regex": f.Search, "$options": "i"}
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("ProblemRepo.List count: %w", mapMongoErr(err))
	}

	var sort bson.D
	switch f.SortBy {
	case "budget":
		sort = bson.D{{Key: "budget_max", Value: -1}}
	case "recent":
		sort = bson.D{{Key: "created_at", Value: -1}}
	default:
		sort = bson.D{{Key: "upvotes", Value: -1}}
	}

	opts := options.Find().SetSort(sort).
		SetSkip(int64((f.Page - 1) * f.Limit)).SetLimit(int64(f.Limit))
	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("ProblemRepo.List find: %w", mapMongoErr(err))
	}
	var out []*domain.Problem
	if err := cur.All(ctx, &out); err != nil {
		return nil, 0, fmt.Errorf("ProblemRepo.List decode: %w", mapMongoErr(err))
	}
	return out, total, nil
}

func (r *ProblemRepo) Update(ctx context.Context, p *domain.Problem) error {
	p.UpdatedAt = time.Now().UTC()
	set := bson.M{
		"title":              p.Title,
		"description":        p.Description,
		"tags":               p.Tags,
		"budget_min":         p.BudgetMin,
		"budget_max":         p.BudgetMax,
		"status":             p.Status,
		"shipped_project_id": p.ShippedProjectID,
		"updated_at":         p.UpdatedAt,
	}
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": p.ID}, bson.M{"$set": set})
	if err != nil {
		return fmt.Errorf("ProblemRepo.Update: %w", mapMongoErr(err))
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("ProblemRepo.Update: %w", domain.ErrNotFound)
	}
	return nil
}

// Claim atomically flips an open problem to claimed. The status guard in the
// filter guarantees only one builder can win the race.
func (r *ProblemRepo) Claim(ctx context.Context, id, builderID string) (*domain.Problem, error) {
	filter := bson.M{"_id": id, "status": domain.ProblemStatusOpen}
	update := bson.M{"$set": bson.M{
		"claimed_by": builderID,
		"status":     domain.ProblemStatusClaimed,
		"updated_at": time.Now().UTC(),
	}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var p domain.Problem
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&p)
	if err == mongo.ErrNoDocuments {
		// Either the problem does not exist or it is no longer open.
		if _, getErr := r.GetByID(ctx, id); getErr != nil {
			return nil, getErr
		}
		return nil, fmt.Errorf("ProblemRepo.Claim: %w", domain.ErrProblemClaimed)
	}
	if err != nil {
		return nil, fmt.Errorf("ProblemRepo.Claim: %w", mapMongoErr(err))
	}
	return &p, nil
}

// AddUpvote atomically adds a unique upvote. $addToSet keeps UpvotedBy unique;
// we only bump the counter when the set actually changed.
func (r *ProblemRepo) AddUpvote(ctx context.Context, id, userID string) (*domain.Problem, error) {
	filter := bson.M{"_id": id, "upvoted_by": bson.M{"$ne": userID}}
	update := bson.M{
		"$addToSet": bson.M{"upvoted_by": userID},
		"$inc":      bson.M{"upvotes": 1},
		"$set":      bson.M{"updated_at": time.Now().UTC()},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var p domain.Problem
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&p)
	if err == mongo.ErrNoDocuments {
		existing, getErr := r.GetByID(ctx, id)
		if getErr != nil {
			return nil, getErr
		}
		// Already upvoted — idempotent success.
		return existing, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ProblemRepo.AddUpvote: %w", mapMongoErr(err))
	}
	return &p, nil
}
