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

func Test_AddOneUser(t *testing.T) {
	fix := struct {
		*fixture
		ownerID string
		userID  string
	}{
		fixture: newFixture(t),
		ownerID: "some_owner_id",
		userID:  "some_user_id",
	}

	groupID, err := fix.app.CreateGroup(fix.ctx, fix.ownerID)
	require.NoError(t, err)

	err = fix.app.AddUserToGroup(fix.ctx, fix.userID, groupID)
	require.NoError(t, err)

	modifiedGroup, err := fix.app.GetGroup(fix.ctx, groupID)
	require.NoError(t, err)

	require.Equal(t, []string{fix.ownerID, fix.userID}, modifiedGroup.Members())
}
