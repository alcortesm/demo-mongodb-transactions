package mongo_test

import (
	"context"
	"testing"
	"time"

	"github.com/alcortesm/demo-mongodb-transactions/internal/domain"
	"github.com/alcortesm/demo-mongodb-transactions/internal/infra/mongo"
	"github.com/alcortesm/demo-mongodb-transactions/internal/testhelp"
	"github.com/stretchr/testify/require"
)

type groupRepoFixture struct {
	// a context with a timeout you can use in your tests
	ctx  context.Context
	repo *mongo.GroupRepo
}

func newGroupRepoFixture(t *testing.T) *groupRepoFixture {
	t.Helper()

	const timeout = 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)

	db := testhelp.NewTestDatabase(t, mongoURI)
	coll := db.Collection("group")
	repo := mongo.NewGroupRepo(coll)

	return &groupRepoFixture{
		ctx:  ctx,
		repo: repo,
	}
}

// Tests you can save a new group and load it later.
func TestGroup_CreateLoad(t *testing.T) {
	fix := struct {
		*groupRepoFixture
		groupID string
		ownerID string
	}{
		groupRepoFixture: newGroupRepoFixture(t),
		groupID:          "group_id",
		ownerID:          "owner_id",
	}

	// GIVEN a group with owned by ownerID in the repo
	group := domain.NewGroup(fix.groupID, fix.ownerID)
	err := fix.repo.Create(fix.ctx, group)
	require.NoError(t, err)

	// WHEN you load the group
	group2, err := fix.repo.Load(fix.ctx, fix.groupID)
	require.NoError(t, err)

	// THEN you get the same group that you saved
	require.Equal(t, group.Snapshot(), group2.Snapshot())
}

// Tests that Update overwrites the document in the db.
func TestGroup_Update(t *testing.T) {
	fix := struct {
		*groupRepoFixture
		groupID string
	}{
		groupRepoFixture: newGroupRepoFixture(t),
		groupID:          "group_id",
	}

	// GIVEN a group owned by fix.onwerID
	group1 := domain.NewGroup(fix.groupID, "owner_id_1")
	group2 := domain.NewGroup(fix.groupID, "owner_id_2")
	require.NotEqual(t, group1.Snapshot(), group2.Snapshot())

	// GIVEN group1 is saved in the db
	err := fix.repo.Create(fix.ctx, group1)
	require.NoError(t, err)

	// WHEN we save group2
	err = fix.repo.Update(fix.ctx, group2)

	// THEN we get no error
	require.NoError(t, err)

	// THEN loading the group id returns group2
	got, err := fix.repo.Load(fix.ctx, fix.groupID)
	require.NoError(t, err)
	require.Equal(t, group2.Snapshot(), got.Snapshot())
}

func TestGroup_LoadNotFound(t *testing.T) {
	fix := struct {
		*groupRepoFixture
	}{
		groupRepoFixture: newGroupRepoFixture(t),
	}

	// WHEN we load a non existing group id
	_, err := fix.repo.Load(fix.ctx, "non_existing_group_id")

	// THEN we get a domain.ErrNotFound error
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrNotFound)
}
