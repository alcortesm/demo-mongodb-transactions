package domain_test

import (
	"testing"

	"github.com/alcortesm/demo-mongodb-transactions/domain"
	"github.com/stretchr/testify/require"
)

func TestGroup_Ctor(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		// GIVEN a user
		user := domain.NewUser("irrelevant_user_id", domain.SkillLevel{})

		// WHEN we create a group with a proper owner
		_, err := domain.NewGroup("irrelevant_group_id", user)

		// THEN we get no error
		require.NoError(t, err)
	})

	t.Run("nil owner", func(t *testing.T) {
		t.Parallel()

		// WHEN we create a group with a nil owner
		_, err := domain.NewGroup("irrelevant_group_id", nil)

		// THEN we get an error
		require.Error(t, err)
	})
}

func TestGroup_ID(t *testing.T) {
	t.Parallel()

	const id = "some_group_id"

	// GIVEN a irrelevant user
	user := domain.NewUser("irrelevant_user_id", domain.SkillLevel{})

	// GIVEN a group with an specific id
	group, err := domain.NewGroup(id, user)
	require.NoError(t, err)

	// WHEN we ask the id of the group
	got := group.ID()

	// THEN we get the correct id
	require.Equal(t, id, got)
}

func TestGroup_OwnerID(t *testing.T) {
	t.Parallel()

	// GIVEN a user
	user := domain.NewUser("some_user_id", domain.SkillLevel{})

	// GIVEN a group owned by the user
	group, err := domain.NewGroup("irrelevant_group_id", user)
	require.NoError(t, err)

	// WHEN we ask for the owner id of the group
	got := group.OwnerID()

	// THEN we get the user id
	require.Equal(t, user.ID(), got)
}

func TestGroup_Members(t *testing.T) {
	t.Parallel()

	// GIVEN user1 and user2
	user1 := domain.NewUser("user_id_1", domain.SkillLevel{})
	user2 := domain.NewUser("user_id_2", domain.SkillLevel{})

	// GIVEN a group with user1 and user2.
	group, err := domain.NewGroup("irrelevant_group_id", user1)
	require.NoError(t, err)
	err = group.AddMember(user2)
	require.NoError(t, err)

	// WHEN we ask for the members
	got := group.Members()

	// THEN we get user1 and user2
	want := []domain.InmutableUser{user1, user2}
	require.ElementsMatch(t, want, got)
}
