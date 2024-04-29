package mongo

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/alcortesm/demo-mongodb-transactions/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
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

// WithTransaction executes callback inside a transaction. If the callback
// returns domain.ErrTransientTransaction it will be retried up to
// maxRetries times.
//
// The callback MUST be idempotent. If maxRetries is 0, the callback won't be
// retried, except for the automatic retry MongoDB does on some very especific
// internal errors.
//
// Errors:
//   - domain.ErrTooManyTransactionRetries if the transaction has failed
//     more than maxRetries times.
//   - whatever non-ErrTransientTransaction errors the callback returns.
func (s *GroupRepo) WithTransaction(
	ctx context.Context,
	callback func(context.Context) error,
	maxRetries uint,
) error {
	session, err := s.transactionSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	sessionCallback := func(txCtx mongo.SessionContext) (interface{}, error) {
		return nil, callback(txCtx)
	}

	for i := range maxRetries {
		if i > 0 {
			log.Printf("retrying transaction, retry %d\n", i)
		}

		switch _, err := session.WithTransaction(ctx, sessionCallback); {
		case err == nil:
			log.Printf("transaction success, attempt: %d\n", i)
			return nil // success
		case errors.Is(err, domain.ErrTransientTransaction):
			log.Printf("transaction failed with transient error, attempt: %d\n", i)
			continue
		default:
			log.Printf("transaction failed with permanent error, attempt: %d", i)
			return err
		}
	}

	return domain.ErrTooManyTransactionRetries
}

// transactionSession creates a session from the client associated with the db
// in s. This session can be used to run transactions on it.
//
// The session is configured with a very conservative set of options that
// allows it to run most kind of transactions in a safe way.
func (r *GroupRepo) transactionSession() (mongo.Session, error) {
	db := r.coll.Database()
	client := db.Client()

	opts := options.Session().
		SetDefaultReadPreference(readpref.Primary()).
		SetDefaultReadConcern(readconcern.Majority()).
		SetDefaultWriteConcern(writeconcern.Majority())

	session, err := client.StartSession(opts)
	if err != nil {
		return nil, fmt.Errorf("starting session: %v", err)
	}

	return session, nil
}
