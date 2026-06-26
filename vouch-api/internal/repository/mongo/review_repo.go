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

// ReviewRepo is the MongoDB implementation of domain.ReviewRepository.
type ReviewRepo struct {
	coll *mongo.Collection
}

// NewReviewRepo constructs a ReviewRepo.
func NewReviewRepo(c *Client) *ReviewRepo {
	return &ReviewRepo{coll: c.DB().Collection("reviews")}
}

func (r *ReviewRepo) Create(ctx context.Context, rv *domain.Review) error {
	if rv.ID == "" {
		rv.ID = primitive.NewObjectID().Hex()
	}
	now := time.Now().UTC()
	rv.CreatedAt, rv.UpdatedAt = now, now
	if _, err := r.coll.InsertOne(ctx, rv); err != nil {
		return fmt.Errorf("ReviewRepo.Create: %w", mapMongoErr(err))
	}
	return nil
}

func (r *ReviewRepo) GetByProjectAndReviewer(ctx context.Context, projectID, reviewerID string) (*domain.Review, error) {
	var rv domain.Review
	err := r.coll.FindOne(ctx, bson.M{"project_id": projectID, "reviewer_id": reviewerID}).Decode(&rv)
	if err != nil {
		return nil, fmt.Errorf("ReviewRepo.GetByProjectAndReviewer: %w", mapMongoErr(err))
	}
	return &rv, nil
}

func (r *ReviewRepo) ListByProject(ctx context.Context, projectID string, page, limit int) ([]*domain.Review, int64, error) {
	filter := bson.M{"project_id": projectID}
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("ReviewRepo.ListByProject count: %w", mapMongoErr(err))
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((page - 1) * limit)).SetLimit(int64(limit))
	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("ReviewRepo.ListByProject find: %w", mapMongoErr(err))
	}
	var out []*domain.Review
	if err := cur.All(ctx, &out); err != nil {
		return nil, 0, fmt.Errorf("ReviewRepo.ListByProject decode: %w", mapMongoErr(err))
	}
	return out, total, nil
}

func (r *ReviewRepo) StatsForProject(ctx context.Context, projectID string) (domain.ReviewStats, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"project_id": projectID}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"count": bson.M{"$sum": 1},
			"avg":   bson.M{"$avg": "$rating"},
		}}},
	}
	cur, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return domain.ReviewStats{}, fmt.Errorf("ReviewRepo.StatsForProject: %w", mapMongoErr(err))
	}
	var rows []struct {
		Count int     `bson:"count"`
		Avg   float64 `bson:"avg"`
	}
	if err := cur.All(ctx, &rows); err != nil {
		return domain.ReviewStats{}, fmt.Errorf("ReviewRepo.StatsForProject decode: %w", mapMongoErr(err))
	}
	if len(rows) == 0 {
		return domain.ReviewStats{}, nil
	}
	return domain.ReviewStats{Count: rows[0].Count, Average: rows[0].Avg}, nil
}
