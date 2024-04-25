package domain_test

import (
	"testing"

	"github.com/alcortesm/demo-mongodb-transactions/domain"
	"github.com/stretchr/testify/require"
)

// Tests that you can recreate a group by taking a snapshot
// and generating the group from it later.
func TestGroupSnapshot_Regenerate(t *testing.T) {
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
}
