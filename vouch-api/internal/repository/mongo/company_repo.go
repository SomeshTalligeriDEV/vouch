package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

const colCompanies = "companies"

// CompanyRepo implements domain.CompanyRepository against MongoDB.
type CompanyRepo struct {
	col *mongo.Collection
}

// NewCompanyRepo constructs a CompanyRepo.
func NewCompanyRepo(c *Client) *CompanyRepo {
	return &CompanyRepo{col: c.db.Collection(colCompanies)}
}

func (r *CompanyRepo) Create(ctx context.Context, c *domain.Company) error {
	if c.ID == "" {
		c.ID = primitive.NewObjectID().Hex()
	}
	c.CreatedAt = time.Now().UTC()
	c.UpdatedAt = c.CreatedAt
	_, err := r.col.InsertOne(ctx, c)
	if err != nil {
		return mapMongoErr(err)
	}
	return nil
}

func (r *CompanyRepo) GetByID(ctx context.Context, id string) (*domain.Company, error) {
	if id == "" {
		return nil, domain.ErrNotFound
	}
	var c domain.Company
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&c); err != nil {
		return nil, mapMongoErr(err)
	}
	return &c, nil
}

func (r *CompanyRepo) GetByEmail(ctx context.Context, email string) (*domain.Company, error) {
	var c domain.Company
	if err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&c); err != nil {
		return nil, mapMongoErr(err)
	}
	return &c, nil
}

func (r *CompanyRepo) GetBySlug(ctx context.Context, slug string) (*domain.Company, error) {
	var c domain.Company
	if err := r.col.FindOne(ctx, bson.M{"slug": slug}).Decode(&c); err != nil {
		return nil, mapMongoErr(err)
	}
	return &c, nil
}

func (r *CompanyRepo) Update(ctx context.Context, c *domain.Company) error {
	oid, err := objectID(c.ID)
	if err != nil {
		return domain.ErrNotFound
	}
	c.UpdatedAt = time.Now().UTC()
	_, err = r.col.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": c})
	return mapMongoErr(err)
}

func (r *CompanyRepo) List(ctx context.Context, page, limit int) ([]*domain.Company, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	skip := int64((page - 1) * limit)
	total, err := r.col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, mapMongoErr(err)
	}
	opts := options.Find().SetSkip(skip).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1})
	cur, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, mapMongoErr(err)
	}
	defer cur.Close(ctx)
	var out []*domain.Company
	if err := cur.All(ctx, &out); err != nil {
		return nil, 0, mapMongoErr(err)
	}
	return out, total, nil
}
