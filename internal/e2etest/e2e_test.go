package e2etest

import (
	"context"
	"testing"
	"time"

	"github.com/alcortesm/demo-mongodb-transactions/internal/application"
	"github.com/alcortesm/demo-mongodb-transactions/internal/infra/mongo"
	"github.com/alcortesm/demo-mongodb-transactions/internal/testhelp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fixture struct {
	// a context with a timeout you can use in your tests
	ctx context.Context
	app *application.App
}

type googleUuider struct{}

func (googleUuider) NewString() string {
	return uuid.NewString()
}

func newFixture(t *testing.T) *fixture {
	t.Helper()

	const timeout = 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)

	db := testhelp.NewTestDatabase(t, mongoURI)
	coll := db.Collection("group")
	groupRepo := mongo.NewGroupRepo(coll)

	app := application.New(googleUuider{}, groupRepo)

	return &fixture{
		ctx: ctx,
		app: app,
	}
}

func Test_CreateGroup(t *testing.T) {
	fix := struct {
		*fixture
		ownerID string
	}{
		fixture: newFixture(t),
		ownerID: "some_owner_id",
	}

	// GIVEN a group owned by fix.ownerID
	groupID, err := fix.app.CreateGroup(fix.ctx, fix.ownerID)
	require.NoError(t, err)

	// WHEN we get the group
	group, err := fix.app.GetGroup(fix.ctx, groupID)
	require.NoError(t, err)

	// THEN the group has the data we expect
	require.Equal(t, groupID, group.ID())
	require.Equal(t, fix.ownerID, group.OwnerID())
	require.Equal(t, []string{fix.ownerID}, group.Members())
}

func Test_AddOneUserToGroup(t *testing.T) {
	fix := struct {
		*fixture
		ownerID string
		userID  string
	}{
		fixture: newFixture(t),
		ownerID: "some_owner_id",
		userID:  "some_user_id",
	}

	// GIVEN a group owned by fix.ownerID
	groupID, err := fix.app.CreateGroup(fix.ctx, fix.ownerID)
	require.NoError(t, err)

	// WHEN we add fix.userID to the group
	err = fix.app.AddUserToGroup(fix.ctx, fix.userID, groupID)
	require.NoError(t, err)

	// THEN the group has the user as a member
	modifiedGroup, err := fix.app.GetGroup(fix.ctx, groupID)
	require.NoError(t, err)
	require.Equal(t, []string{fix.ownerID, fix.userID}, modifiedGroup.Members())
}
