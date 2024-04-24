package domain

// User represent an authenticated user in our system.
//
// Users are not zero-value-safe.
type User struct {
	id         string
	skillLevel SkillLevel
}

// InmutableUser is an inmutable view of a User.
type InmutableUser interface {
	ID() string
	SkillLevel() SkillLevel
}

// User returns a new user.
func NewUser(id string, skillLevel SkillLevel) *User {
	return &User{
		id:         id,
		skillLevel: skillLevel,
	}
}

// ID returns the user id.
func (u *User) ID() string {
	return u.id
}

// SkillLevel returns the user skill level.
func (u *User) SkillLevel() SkillLevel {
	return u.skillLevel
}
