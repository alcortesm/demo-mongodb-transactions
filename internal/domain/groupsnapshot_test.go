package domain_test

import (
	"fmt"
	"testing"

	"github.com/alcortesm/demo-mongodb-transactions/internal/domain"
	"github.com/stretchr/testify/require"
)

// Tests that you can recreate a group by taking a snapshot
// and generating the group from it later.
func TestGroupSnapshot_Regenerate(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		const (
			user1 = "user_1_id"
			user2 = "user_2_id"
		)

		// GIVEN a group with user1 and user2 as members
		group := domain.NewGroup("irrelevant_group_id", user1)
		err := group.AddMember(user2)
		require.NoError(t, err)

		// GIVEN a snapshot of the group
		snapshot := group.Snapshot()

		// WHEN you recreate the group from the snapshot
		group2, err := snapshot.Regenerate()
		require.NoError(t, err)

		// THEN the regenerated group looks the same as the original group
		require.Equal(t, group.ID(), group2.ID())
		require.Equal(t, group.OwnerID(), group2.OwnerID())
		require.Equal(t, group.Members(), group2.Members())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		tooManyMembers := func(t *testing.T) *domain.GroupSnapshot {
			t.Helper()

			g := &domain.GroupSnapshot{
				ID:      "irrelevant_group_id",
				OwnerID: "irrelevant_owner_id",
				Members: []string{"irrelevant_owner_id"},
			}

			for i := range domain.MaxMembers {
				g.Members = append(g.Members, fmt.Sprintf("user_id_%d", i))
			}

			require.True(t, len(g.Members) > domain.MaxMembers)

			return g
		}

		subtests := []struct {
			name         string
			snapshot     *domain.GroupSnapshot
			errorContent string
		}{
			{
				name: "empty id",
				snapshot: &domain.GroupSnapshot{
					ID:      "",
					OwnerID: "irrelevant_owner_id",
					Members: []string{"user_id_1", "user_id_2"},
				},
				errorContent: "empty id",
			},
			{
				name:         "too many members",
				snapshot:     tooManyMembers(t),
				errorContent: "too many members",
			},
			{
				name: "owner is not a member",
				snapshot: &domain.GroupSnapshot{
					ID:      "irrelevant_group_id",
					OwnerID: "a",
					Members: []string{"b", "c"},
				},
				errorContent: "not member",
			},
			{
				name: "no members",
				snapshot: &domain.GroupSnapshot{
					ID:      "irrelevant_group_id",
					OwnerID: "irrelevant_owner_id",
					Members: []string{},
				},
				errorContent: "empty members",
			},
			{
				name: "nil members",
				snapshot: &domain.GroupSnapshot{
					ID:      "irrelevant_group_id",
					OwnerID: "irrelevant_owner_id",
					Members: nil,
				},
				errorContent: "empty members",
			},
			{
				name: "empty owner",
				snapshot: &domain.GroupSnapshot{
					ID:      "irrelevant_group_id",
					OwnerID: "",
					Members: []string{"user_id_1", "user_id_2"},
				},
				errorContent: "empty owner id",
			},
		}

		for _, test := range subtests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				// WHEN we regenerate the snapshot
				_, err := test.snapshot.Regenerate()

				// THEN we must get an error with the content we want
				require.Error(t, err)
				require.ErrorContains(t, err, test.errorContent)
			})
		}
	})
}
