package mongo

import (
	"context"
	"fmt"

	"github.com/alcortesm/demo-mongodb-transactions/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupRepo struct {
	coll *mongo.Collection
}

func NewGroupRepo(coll *mongo.Collection) *GroupRepo {
	return &GroupRepo{
		coll: coll,
	}
}

// Create stores the group in the database as a new document.
//
// Error:
//   - domain.ErrTransientTransaction if the operation failed during a transaction
//     that can be retried.
func (r *GroupRepo) Create(ctx context.Context, group *domain.Group) error {
	doc := newGroupDoc(group)

	_, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("replacing: %w", domainError(err))
	}

	return nil
}

// Update overwrites the group document in the database.
//
// Error:
//   - domain.ErrNotFound if the group is not found
//   - domain.ErrTransientTransaction if the operation failed during a transaction
//     that can be retried.
func (r *GroupRepo) Update(ctx context.Context, group *domain.Group) error {
	filter := bson.M{
		"_id": group.ID(),
	}

	doc := newGroupDoc(group)

	result, err := r.coll.ReplaceOne(ctx, filter, doc)
	if err != nil {
		return fmt.Errorf("replacing: %w", domainError(err))
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
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
