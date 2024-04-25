package domain_test

import (
	"fmt"
	"testing"

	"github.com/alcortesm/demo-mongodb-transactions/domain"
	"github.com/stretchr/testify/require"
)

func TestGroup_ID(t *testing.T) {
	t.Parallel()

	const id = "some_group_id"

	// GIVEN a group with an specific id
	group := domain.NewGroup(id, "irrelevant_owner_id")

	// WHEN we ask the id of the group
	got := group.ID()

	// THEN we get the correct id
	require.Equal(t, id, got)
}

func TestGroup_OwnerID(t *testing.T) {
	t.Parallel()

	const ownerID = "owner_id"

	// GIVEN a group owned by the user
	group := domain.NewGroup("irrelevant_group_id", ownerID)

	// WHEN we ask for the owner id of the group
	got := group.OwnerID()

	// THEN we get the user id
	require.Equal(t, ownerID, got)
}

func TestGroup_AddMembers(t *testing.T) {
	t.Parallel()

	t.Run("add less than MaxMembers", func(t *testing.T) {
		const (
			ownerID = "owner_id"
			user1ID = "user_id_1"
			user2ID = "user_id_2"
			user3ID = "user_id_3"
		)

		subtests := []struct {
			name        string
			usersToAdd  []string
			wantMembers []string
		}{
			{
				name:        "only owner",
				usersToAdd:  nil,
				wantMembers: []string{ownerID},
			},
			{
				name:        "owner and user1",
				usersToAdd:  []string{user1ID},
				wantMembers: []string{ownerID, user1ID},
			},
			{
				name:        "owner user1 and user2",
				usersToAdd:  []string{user1ID, user2ID},
				wantMembers: []string{ownerID, user1ID, user2ID},
			},
			{
				name:        "owner user1 user2 and user3",
				usersToAdd:  []string{user1ID, user2ID, user3ID},
				wantMembers: []string{ownerID, user1ID, user2ID, user3ID},
			},
		}

		for _, test := range subtests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				// GIVEN a group owned by ownerID
				group := domain.NewGroup("irrelevant_group_id", ownerID)

				// WHEN we add test.usersToAdd to the group
				for _, id := range test.usersToAdd {
					err := group.AddMember(id)
					require.NoErrorf(t, err, "adding user %s to group", id)
				}

				// THEN the members are the ones we want
				{
					got := group.Members()
					require.Equal(t, test.wantMembers, got)
				}

				// THEN the members count is the one we want
				{
					got := group.NumMembers()
					require.Equal(t, len(test.wantMembers), got)
				}

			})
		}
	})

	t.Run("add already existing member", func(t *testing.T) {
		const (
			ownerID = "owner_id" // owner of the group
		)

		// GIVEN a group owned by ownerID
		group := domain.NewGroup("irrelevant_group_id", ownerID)
		membersBefore := group.Members()

		// WHEN we try to add an already existing member
		err := group.AddMember(ownerID)

		// THEN we get success
		require.NoError(t, err)

		// THEN there is no change in the list of members
		{
			got := group.Members()
			require.Equal(t, membersBefore, got)
		}
	})

	t.Run("group full", func(t *testing.T) {
		// userID builds user id in the form: "user_id_<n>"
		userID := func(n int) string {
			return fmt.Sprintf("user_id_%d", n)
		}

		// GIVEN a full group
		group := domain.NewGroup("irrelevant_group_id", userID(0))
		for i := 1; i < domain.MaxMembers; i++ {
			err := group.AddMember(userID(i))
			require.NoErrorf(t, err, "adding user with index #%d", i)
		}
		require.Equal(t, domain.MaxMembers, group.NumMembers())

		// WHEN we try to add one more user
		err := group.AddMember("one_more_user_id")

		// THEN we get ErrGroupFull
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrGroupFull)
	})
}
