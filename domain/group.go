package domain

import (
	"errors"
)

// Group represents a group of users.
//
// Invariants:
//   - must have at least one member.
//   - cannot have more than MaxMembers members.
//   - must have an owner, which is one of its members.
type Group struct {
	id      string
	ownerID string
	members map[string]InmutableUser
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
func NewGroup(id string, owner *User) (*Group, error) {
	if owner == nil {
		return nil, errors.New("nil owner")
	}

	return &Group{
		id:      id,
		ownerID: owner.ID(),
		members: map[string]InmutableUser{
			owner.ID(): owner,
		},
	}, nil
}

// ID returns the group id.
func (g *Group) ID() string {
	return g.id
}

// OwnerID returns the onwer id.
func (g *Group) OwnerID() string {
	return g.ownerID
}

// Members returns a slice with the members in a non-specified order.
func (g *Group) Members() []InmutableUser {
	result := make([]InmutableUser, 0, len(g.members))

	for _, v := range g.members {
		result = append(result, v)
	}

	return result
}

// AddMember adds a user to the group.
//
// If the user was already a member, it is no-op and returns nil.
//
// Returns:
// - ErrMaxMembers if the group is already full
func (g *Group) AddMember(user *User) error {
	if len(g.members) >= MaxMembers {
		return ErrGroupFull
	}

	g.members[user.ID()] = user

	return nil
}
