package domain

import (
	"fmt"
)

// SkillLevel represents how skillful a user is.
//
// The skill level of a user can be:
//
//   - unknown (zero value)
//   - a value in the [MinSkillLevel, MaxSkillLevel] rage (see [NewSkillLevel])
//
// SkillLevel are inmutable values and are comparable with '=='.
//
// The zero value represents an unknown skill level for a user.
type SkillLevel struct {
	known bool
	value int
}

const (
	MinSkillLevel = 0
	MaxSkillLevel = 10
)

// NewSkillLevel returns a new skill level of value n.
func NewSkillLevel(n int) (SkillLevel, error) {
	if n < MinSkillLevel {
		return SkillLevel{}, fmt.Errorf("invalid skill level %d, must be ≥ %d", n, MinSkillLevel)
	}

	if n > MaxSkillLevel {
		return SkillLevel{}, fmt.Errorf("invalid skill level %d, must be ≤ %d", n, MaxSkillLevel)
	}

	return SkillLevel{
		known: true,
		value: n,
	}, nil
}

// Value returns the skill level and true if it is known, otherwise it returns 0 and false.
func (s SkillLevel) Value() (int, bool) {
	if s.IsKnown() {
		return s.value, true
	}

	return 0, false
}

// IsKnown returns if the skill level is known.
func (s SkillLevel) IsKnown() bool {
	return s.known
}

// String returns a human readable representation of the skill level for
// debugging purposes.
func (s SkillLevel) String() string {
	value, ok := s.Value()
	if !ok {
		return "unknown_skill_level"
	}

	return fmt.Sprintf("%d", value)
}
