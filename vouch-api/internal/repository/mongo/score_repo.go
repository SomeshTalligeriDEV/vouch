package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// ScoreRepo is the MongoDB implementation of domain.ScoreRepository.
type ScoreRepo struct {
	coll *mongo.Collection
}

// NewScoreRepo constructs a ScoreRepo.
func NewScoreRepo(c *Client) *ScoreRepo {
	return &ScoreRepo{coll: c.DB().Collection("builder_scores")}
}

func (r *ScoreRepo) GetByBuilderID(ctx context.Context, builderID string) (*domain.BuilderScore, error) {
	var s domain.BuilderScore
	if err := r.coll.FindOne(ctx, bson.M{"builder_id": builderID}).Decode(&s); err != nil {
		return nil, fmt.Errorf("ScoreRepo.GetByBuilderID: %w", mapMongoErr(err))
	}
	return &s, nil
}

func (r *ScoreRepo) Upsert(ctx context.Context, s *domain.BuilderScore) error {
	s.UpdatedAt = time.Now().UTC()
	set := bson.M{
		"builder_id":        s.BuilderID,
		"total_score":       s.TotalScore,
		"tier":              s.Tier,
		"breakdown":         s.Breakdown,
		"stripe_verified":   s.StripeVerified,
		"stripe_multiplier": s.StripeMultiplier,
		"calculated_at":     s.CalculatedAt,
		"updated_at":        s.UpdatedAt,
	}
	_, err := r.coll.UpdateOne(ctx,
		bson.M{"builder_id": s.BuilderID},
		bson.M{"$set": set},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("ScoreRepo.Upsert: %w", mapMongoErr(err))
	}
	return nil
}

func (r *ScoreRepo) TopBuilders(ctx context.Context, limit int) ([]*domain.BuilderScore, error) {
	opts := options.Find().SetSort(bson.D{{Key: "total_score", Value: -1}}).SetLimit(int64(limit))
	cur, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("ScoreRepo.TopBuilders: %w", mapMongoErr(err))
	}
	var out []*domain.BuilderScore
	if err := cur.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("ScoreRepo.TopBuilders decode: %w", mapMongoErr(err))
	}
	return out, nil
}
