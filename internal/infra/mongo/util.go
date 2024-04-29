package mongo

import (
	"errors"
	"fmt"

	"github.com/alcortesm/demo-mongodb-transactions/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

// domainError translates between MongoDB errors and domain errors, hiding
// MongoDB implementation details and returning well known domain errors:
//
//   - domain.ErrTransientTransaction: on mongo errors with the
//     driver.TransientTransactionError label
//
//   - domain.ErrNotFound: whatever you were looking for, it has not been found.
func domainError(err error) error {
	// check for transient transaction errors.
	for current := err; current != nil; current = errors.Unwrap(current) {
		if le, ok := current.(mongo.LabeledError); ok && le.HasErrorLabel(driver.TransientTransactionError) {
			return fmt.Errorf("%w: %v", domain.ErrTransientTransaction, err)
		}
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.ErrNotFound
	}

	return fmt.Errorf("%v", err) // do not expose any mongoDB errors
}
