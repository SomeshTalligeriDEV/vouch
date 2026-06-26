package mongo

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SomeshTalligeriDEV/vouch-api/internal/domain"
)

// bsonDoc builds a single-key bson.D for index definitions.
func bsonDoc(key string, order int) bson.D {
	return bson.D{{Key: key, Value: order}}
}

// compoundKey builds a two-key ascending bson.D for compound indexes.
func compoundKey(a, b string) bson.D {
	return bson.D{{Key: a, Value: 1}, {Key: b, Value: 1}}
}

// objectID parses a hex string into an ObjectID, mapping bad input to
// domain.ErrNotFound (a malformed ID can never match a real document).
func objectID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, domain.ErrNotFound
	}
	return oid, nil
}

// mapMongoErr normalizes mongo errors into domain errors.
func mapMongoErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.ErrNotFound
	}
	if mongo.IsDuplicateKeyError(err) {
		return domain.ErrAlreadyExists
	}
	return err
}
