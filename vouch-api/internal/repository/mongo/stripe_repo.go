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

// StripeRepo is the MongoDB implementation of domain.StripeRepository.
type StripeRepo struct {
	coll *mongo.Collection
}

// NewStripeRepo constructs a StripeRepo.
func NewStripeRepo(c *Client) *StripeRepo {
	return &StripeRepo{coll: c.DB().Collection("stripe_snapshots")}
}

func (r *StripeRepo) Save(ctx context.Context, s *domain.StripeSnapshot) error {
	if s.ID == "" {
		s.ID = primitive.NewObjectID().Hex()
	}
	if s.VerifiedAt.IsZero() {
		s.VerifiedAt = time.Now().UTC()
	}
	if _, err := r.coll.InsertOne(ctx, s); err != nil {
		return fmt.Errorf("StripeRepo.Save: %w", mapMongoErr(err))
	}
	return nil
}

func (r *StripeRepo) Latest(ctx context.Context, builderID string) (*domain.StripeSnapshot, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "verified_at", Value: -1}})
	var s domain.StripeSnapshot
	if err := r.coll.FindOne(ctx, bson.M{"builder_id": builderID}, opts).Decode(&s); err != nil {
		return nil, fmt.Errorf("StripeRepo.Latest: %w", mapMongoErr(err))
	}
	return &s, nil
}
