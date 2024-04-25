package domain

import (
	"sort"
)

type empty struct{}

// Group represents a group of users.
//
// Invariants:
//   - must have at least one member.
//   - cannot have more than MaxMembers members.
//   - must have an owner, which is one of its members.
type Group struct {
	id      string
	ownerID string
	members map[string]empty
}

// MaxMembers is maximum number of members in a group.
//
// You can safely increase this value at any time, but if you want to reduce
// it, you will need to ensure no groups in the database have more members than
// the new value.
//
// MaxMembers must be greater than 0, so groups have at least one member.
const MaxMembers = 5

// NewGroup creates a new group owned by owner.
//
// The owner is required.
func NewGroup(id string, ownerID string) *Group {
	return &Group{
		id:      id,
		ownerID: ownerID,
		members: map[string]empty{
			ownerID: empty{},
		},
	}
}

// ID returns the group id.
func (g *Group) ID() string {
	return g.id
}

// OwnerID returns the onwer id.
func (g *Group) OwnerID() string {
	return g.ownerID
}

// AddMember adds a user to the group.
//
// If the user was already a member, it is no-op and returns nil.
//
// Returns:
// - ErrMaxMembers if the group is already full
func (g *Group) AddMember(id string) error {
	if len(g.members) >= MaxMembers {
		return ErrGroupFull
	}

	g.members[id] = empty{}

	return nil
}

// Members returns a slice with the members id sorted alphabetically.
func (g *Group) Members() []string {
	result := make([]string, 0, len(g.members))

	for id := range g.members {
		result = append(result, id)
	}

	sort.Strings(result)

	return result
}

// NumMembers returns the count of members in the group.
func (g *Group) NumMembers() int {
	return len(g.members)
}

// Snapshot returns a snapshot of the internal state of the group.
func (g *Group) Snapshot() *GroupSnapshot {
	return &GroupSnapshot{
		ID:      g.ID(),
		OwnerID: g.OwnerID(),
		Members: g.Members(),
	}
}
