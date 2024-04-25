package domain

import (
	"errors"
	"fmt"
	"slices"
)

// GroupSnapshot represent the internal state of a group.
//
// It can be used to persists a group, by saving the snapshot and recreating
// the group from it later.
type GroupSnapshot struct {
	ID      string
	OwnerID string
	// IDs of the members in alphabetical order.
	Members []string
}

// Regenerate creates a group from the internal state represented by s.
//
// Returns an error if the internal state represented by s would lead to an
// invalid group.
func (s *GroupSnapshot) Regenerate() (*Group, error) {
	if s.ID == "" {
		return nil, errors.New("empty id")
	}

	if len(s.Members) > MaxMembers {
		return nil, fmt.Errorf("too many members (%d)", len(s.Members))
	}

	if !slices.Contains(s.Members, s.OwnerID) {
		return nil, fmt.Errorf("owner (%s) is not member (%s)", s.OwnerID, s.Members)
	}

	g := &Group{
		id:      s.ID,
		ownerID: s.OwnerID,
		members: map[string]empty{},
	}

	for _, id := range s.Members {
		g.members[id] = empty{}
	}

	return g, nil
}
