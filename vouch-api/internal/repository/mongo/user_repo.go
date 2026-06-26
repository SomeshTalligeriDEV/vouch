package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// UserRepo is the MongoDB implementation of domain.UserRepository.
type UserRepo struct {
	coll *mongo.Collection
}

// NewUserRepo constructs a UserRepo.
func NewUserRepo(c *Client) *UserRepo {
	return &UserRepo{coll: c.DB().Collection("users")}
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	if u.ID == "" {
		u.ID = primitive.NewObjectID().Hex()
	}
	now := time.Now().UTC()
	u.CreatedAt, u.UpdatedAt = now, now
	if _, err := r.coll.InsertOne(ctx, u); err != nil {
		return fmt.Errorf("UserRepo.Create: %w", mapMongoErr(err))
	}
	return nil
}

func (r *UserRepo) one(ctx context.Context, filter bson.M, op string) (*domain.User, error) {
	var u domain.User
	if err := r.coll.FindOne(ctx, filter).Decode(&u); err != nil {
		return nil, fmt.Errorf("UserRepo.%s: %w", op, mapMongoErr(err))
	}
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return r.one(ctx, bson.M{"_id": id}, "GetByID")
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return r.one(ctx, bson.M{"username": username}, "GetByUsername")
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.one(ctx, bson.M{"email": email}, "GetByEmail")
}

func (r *UserRepo) GetByGitHubID(ctx context.Context, githubID int64) (*domain.User, error) {
	return r.one(ctx, bson.M{"github_id": githubID}, "GetByGitHubID")
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	u.UpdatedAt = time.Now().UTC()
	set := bson.M{
		"name":           u.Name,
		"bio":            u.Bio,
		"avatar_url":     u.AvatarURL,
		"website_url":    u.WebsiteURL,
		"twitter_handle": u.TwitterHandle,
		"role":           u.Role,
		"is_verified":    u.IsVerified,
		"updated_at":     u.UpdatedAt,
	}
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{"$set": set})
	if err != nil {
		return fmt.Errorf("UserRepo.Update: %w", mapMongoErr(err))
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("UserRepo.Update: %w", domain.ErrNotFound)
	}
	return nil
}

func (r *UserRepo) SetStripeAccount(ctx context.Context, id, stripeAccountID string) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id},
		bson.M{"$set": bson.M{"stripe_account_id": stripeAccountID, "updated_at": time.Now().UTC()}})
	if err != nil {
		return fmt.Errorf("UserRepo.SetStripeAccount: %w", mapMongoErr(err))
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("UserRepo.SetStripeAccount: %w", domain.ErrNotFound)
	}
	return nil
}
