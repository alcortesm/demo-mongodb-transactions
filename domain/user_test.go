package domain_test

import (
	"testing"

	"github.com/alcortesm/demo-mongodb-transactions/domain"
	"github.com/alcortesm/demo-mongodb-transactions/testhelp"
	"github.com/stretchr/testify/require"
)

func TestUser_ID(t *testing.T) {
	t.Parallel()

	const userID = "some_user_id"

	// GIVEN an irrelevant skill level for the user
	skillLevel := domain.SkillLevel{}

	// GIVEN a user with a given id
	user := domain.NewUser(userID, skillLevel)

	// WHEN we ask for the user id
	got := user.ID()

	// THEN we get the correct user id
	require.Equal(t, userID, got)
}

func TestUser_SkillLevel(t *testing.T) {
	t.Parallel()

	skillLevels := []domain.SkillLevel{
		domain.SkillLevel{},
		testhelp.SkillLevel(t, domain.MinSkillLevel),
		testhelp.SkillLevel(t, domain.MaxSkillLevel),
	}

	for _, sl := range skillLevels {
		t.Run(sl.String(), func(t *testing.T) {
			t.Parallel()

			// GIVEN a user with the given skill level
			user := domain.NewUser("irrelevant_user_id", sl)

			// WHEN we ask for their skill level
			got := user.SkillLevel()

			// THEN we get the correct skill level
			require.Equal(t, sl, got)
		})
	}
}
