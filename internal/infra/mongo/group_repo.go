package mongo

import (
	"context"
	"fmt"

	"github.com/alcortesm/demo-mongodb-transactions/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupRepo struct {
	coll *mongo.Collection
}

func NewGroupRepo(coll *mongo.Collection) *GroupRepo {
	return &GroupRepo{
		coll: coll,
	}
}

// Save persists the group into the database, overwriting any document with the same id.
//
// Error:
//   - domain.ErrTransientTransaction if the operation failed during a transaction
//     that can be retried.
func (r *GroupRepo) Save(ctx context.Context, group *domain.Group) error {
	filter := bson.M{
		"_id": group.ID(),
	}

	doc := newGroupDoc(group)

	opts := options.Replace().SetUpsert(true)

	_, err := r.coll.ReplaceOne(ctx, filter, doc, opts)
	if err != nil {
		return fmt.Errorf("replacing: %w", domainError(err))
	}

	return nil
}

// Load returns a group with the give id from the database.
//
// Errors:
//   - domain.ErrNotFound if there is no group with the given ID
//   - domain.ErrTransientTransaction if the operation failed during a transaction
//     that can be retried.
func (r *GroupRepo) Load(ctx context.Context, id string) (*domain.Group, error) {
	filter := bson.M{
		"_id": id,
	}

	doc := new(groupDoc)

	err := r.coll.FindOne(ctx, filter).Decode(doc)
	if err != nil {
		return nil, domainError(err)
	}

	return doc.group()
}
