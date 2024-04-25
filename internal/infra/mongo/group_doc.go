package mongo

import "github.com/alcortesm/demo-mongodb-transactions/internal/domain"

// groupDoc is a Mongo document representing a group
type groupDoc struct {
	ID      string   `bson:"_id"`
	OwnerID string   `bson:"owner_id"`
	Members []string `bson:"members"`
}

func newGroupDoc(group *domain.Group) *groupDoc {
	s := group.Snapshot()

	doc := &groupDoc{
		ID:      s.ID,
		OwnerID: s.OwnerID,
		Members: s.Members,
	}

	return doc
}

// group returns the domain.Group represented by docGroup.
func (d *groupDoc) group() (*domain.Group, error) {
	s := domain.GroupSnapshot(*d)
	return s.Regenerate()
}
