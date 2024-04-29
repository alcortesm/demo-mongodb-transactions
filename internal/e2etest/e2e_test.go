package e2etest

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alcortesm/demo-mongodb-transactions/internal/application"
	"github.com/alcortesm/demo-mongodb-transactions/internal/domain"
	"github.com/alcortesm/demo-mongodb-transactions/internal/infra/mongo"
	"github.com/alcortesm/demo-mongodb-transactions/internal/testhelp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

	const timeout = 10 * time.Second
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

// Test the app layer respects the Group invariants while adding users:
//
// Let's make many AddUserToGroup requests concurrently, more than the maximum
// number of members allowed in a group.
//
// If the app layers respects the Group invariants, only MaxMembers-1 requests
// will be successful, the rest will get a domain.ErrFullGroup error.
//
// For example, if MaxMembers is 5 and we sent 10 requests to add users:
//   - 4 will be added correctly (5 minus the owner, which was already a member)
//   - 6 will fail receive a domain.ErrFullGroup error
func Test_Concurrency_AddLotsOfUsersConcurrentlyToGroup(t *testing.T) {
	fix := struct {
		*fixture
		ownerID   string
		userCount int
	}{
		fixture:   newFixture(t),
		ownerID:   "some_owner_id",
		userCount: 10,
	}

	// make sure the are trying to add more than the maximum number of users
	// allowed in a group
	require.GreaterOrEqual(t, fix.userCount, domain.MaxMembers)

	// we will run the test twice: first without transactions in the app, then
	// with transactions enabled
	subtests := []struct {
		name               string
		options            []application.Option
		wantSuccessCount   int
		wantFullGroupCount int
	}{
		{
			// when adding users concurrently with transactions disabled all
			// AddUserToGroup requests will be successful and we will never get
			// an ErrFullGroup error
			name:               "transactions disabled",
			options:            nil, // transactions are disabled by default
			wantSuccessCount:   fix.userCount,
			wantFullGroupCount: 0,
		},
		{
			// when adding users concurrently with transactions ENABLED we will
			// only be able to add a few users, until the group is full, then
			// we will get a bunch of ErrFullGroup errors for the rest of the
			// requests
			name:               "transactions enabled",
			options:            []application.Option{application.EnableTransactions{}},
			wantSuccessCount:   domain.MaxMembers - 1,                   // the owner already counts as a member
			wantFullGroupCount: fix.userCount - (domain.MaxMembers - 1), // the remaininig requests
		},
	}

	for _, test := range subtests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			// userID returns the id of a user based on the number n, for example, "user_id_0042"
			userID := func(n int) string { return fmt.Sprintf("user_id_%02d", n) }

			// GIVEN a group owned by fix.ownerID
			groupID, err := fix.app.CreateGroup(fix.ctx, fix.ownerID)
			require.NoError(t, err)

			// GIVEN more users than we can fit in the group
			users := make([]string, 0, fix.userCount)
			for i := range fix.userCount {
				users = append(users, userID(i))
			}

			// WHEN we add all the users to the group at the same time and keep track
			// of how many requests got success vs how many requests got a
			// domain.ErrGroupFull error.
			var successCount, fullGroupCount int
			{
				var wg sync.WaitGroup
				wg.Add(len(users))

				results := make([]error, len(users))
				for i, id := range users {
					go func() {
						defer wg.Done()

						option := append(
							// introduce an artificial delay in the app layer to improve the chance of
							// processing all the requests at the same time
							[]application.Option{application.DelayBeforeUpdating(500 * time.Millisecond)},
							test.options...,
						)

						results[i] = fix.app.AddUserToGroup(
							fix.ctx,
							id,
							groupID,
							option...,
						)
					}()
				}

				wg.Wait()
				for i, err := range results {
					switch {
					case err == nil:
						successCount++
					case errors.Is(err, domain.ErrGroupFull):
						fullGroupCount++
					default:
						t.Errorf("adding %s: %v", userID(i), err)
					}
				}
			}

			// THEN the number of requests that got a successful reponse is test.wantSuccessCount
			// and the number of requests that got an ErrGroupFull error is test.wantFullGroupCount
			assert.Equal(t, test.wantSuccessCount, successCount, "wrong count of successful calls")
			assert.Equal(t, test.wantFullGroupCount, fullGroupCount, "wrong count of ErrGroupFull received")

		})
	}
}
