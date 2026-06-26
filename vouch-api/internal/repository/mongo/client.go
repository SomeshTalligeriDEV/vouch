package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client wraps a connected MongoDB database handle.
type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

// Connect dials MongoDB and pings to verify connectivity.
func Connect(ctx context.Context, uri, dbName string) (*Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cl, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo.Connect: %w", err)
	}
	if err := cl.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo.Connect ping: %w", err)
	}
	return &Client{client: cl, db: cl.Database(dbName)}, nil
}

// DB returns the underlying database handle.
func (c *Client) DB() *mongo.Database { return c.db }

// Disconnect cleanly closes the connection.
func (c *Client) Disconnect(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// EnsureIndexes creates the indexes Vouch relies on. Safe to call repeatedly.
func (c *Client) EnsureIndexes(ctx context.Context) error {
	uniq := options.Index().SetUnique(true)

	specs := []struct {
		coll  string
		model mongo.IndexModel
	}{
		{"users", mongo.IndexModel{Keys: bsonDoc("email", 1), Options: uniq}},
		{"users", mongo.IndexModel{Keys: bsonDoc("username", 1), Options: uniq}},
		{"users", mongo.IndexModel{Keys: bsonDoc("github_id", 1), Options: uniq}},
		{"projects", mongo.IndexModel{Keys: bsonDoc("slug", 1), Options: uniq}},
		{"projects", mongo.IndexModel{Keys: bsonDoc("builder_id", 1)}},
		{"projects", mongo.IndexModel{Keys: bsonDoc("status", 1)}},
		{"builder_scores", mongo.IndexModel{Keys: bsonDoc("builder_id", 1), Options: uniq}},
		{"builder_scores", mongo.IndexModel{Keys: bsonDoc("total_score", -1)}},
		{"problems", mongo.IndexModel{Keys: bsonDoc("slug", 1), Options: uniq}},
		{"problems", mongo.IndexModel{Keys: bsonDoc("status", 1)}},
		{"reviews", mongo.IndexModel{Keys: bsonDoc("project_id", 1)}},
		{"reviews", mongo.IndexModel{Keys: compoundKey("project_id", "reviewer_id"), Options: uniq}},
		{"stripe_snapshots", mongo.IndexModel{Keys: bsonDoc("builder_id", 1)}},
	}

	for _, s := range specs {
		if _, err := c.db.Collection(s.coll).Indexes().CreateOne(ctx, s.model); err != nil {
			return fmt.Errorf("mongo.EnsureIndexes %s: %w", s.coll, err)
		}
	}
	return nil
}
