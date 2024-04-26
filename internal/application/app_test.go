package application_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/alcortesm/demo-mongodb-transactions/internal/application"
	"github.com/alcortesm/demo-mongodb-transactions/internal/domain"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type fixture struct {
	app    *application.App
	uuider *MockUuider
	store  *MockStore
}

func newFixture(t *testing.T) *fixture {
	ctrl := gomock.NewController(t)
	uuider := NewMockUuider(ctrl)
	store := NewMockStore(ctrl)

	app := application.New(uuider, store)

	return &fixture{
		app:    app,
		uuider: uuider,
		store:  store,
	}
}

func TestCreateGroup(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		fix := struct {
			*fixture
			ownerID string
			groupID string
		}{
			fixture: newFixture(t),
			ownerID: "some_owner_id",
			groupID: "some_group_id",
		}

		// GIVEN a uuider that returns the new group id
		fix.uuider.EXPECT().
			NewString().
			Return(fix.groupID)

		// GIVEN a groupRepo that fails the test if passed the wrong group
		{
			want := domain.NewGroup(fix.groupID, fix.ownerID)
			fix.store.EXPECT().
				Create(gomock.Any(), want).
				Return(nil)
		}

		// WHEN we create a group
		id, err := fix.app.CreateGroup(context.Background(), fix.ownerID)
		require.NoError(t, err)

		// THEN we return the new group id
		require.Equal(t, fix.groupID, id)
	})

	t.Run("groupRepo create error", func(t *testing.T) {
		fix := struct {
			*fixture
		}{
			fixture: newFixture(t),
		}

		// GIVEN a uuider that returns the same new group id every time
		fix.uuider.EXPECT().
			NewString().
			Return("irrelevant_group_id")

		// GIVEN a groupRepo that fails to create
		cause := errors.New("some_repo_error")
		fix.store.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(cause)

		// WHEN we create a group
		_, err := fix.app.CreateGroup(context.Background(), "irrelevant_owner_id")

		// THEN we get the error we expect
		require.Error(t, err)
		require.ErrorContains(t, err, cause.Error())
	})
}

func TestGetGroup(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		fix := struct {
			*fixture
			groupID string
		}{
			fixture: newFixture(t),
			groupID: "some_group_id",
		}

		// GIVEN a groupRepo that fails the test if passed the wrong group id
		// and that returns a group
		group := domain.NewGroup(fix.groupID, "irrelevant_owner_id")
		fix.store.EXPECT().
			Load(gomock.Any(), fix.groupID).
			Return(group, nil)

		// WHEN we get the group
		got, err := fix.app.GetGroup(context.Background(), fix.groupID)
		require.NoError(t, err)

		// THEN we return the group we want
		require.Equal(t, group.Snapshot(), got.Snapshot())
	})

	t.Run("groupRepo load error", func(t *testing.T) {
		fix := struct {
			*fixture
		}{
			fixture: newFixture(t),
		}

		// GIVEN a groupRepo that fails to get
		cause := errors.New("some_repo_error")
		fix.store.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil, cause)

		// WHEN we get a group
		_, err := fix.app.GetGroup(context.Background(), "irrelevant_group_id")

		// THEN we get the error we expect
		require.Error(t, err)
		require.ErrorContains(t, err, cause.Error())
	})
}

func TestAddUserToGroup(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		fix := struct {
			*fixture
			groupID string
			userID  string
		}{
			fixture: newFixture(t),
			groupID: "some_group_id",
			userID:  "some_user_id",
		}

		// GIVEN a groupRepo expecting a Load with for the right group
		group := domain.NewGroup(fix.groupID, "irrelevant_owner_id")
		fix.store.EXPECT().
			Load(gomock.Any(), fix.groupID).
			Return(group, nil)

		// GIVEN-THEN a groupRepo expecting an Update with the new user as a member
		fix.store.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			DoAndReturn(func(_ any, got *domain.Group) error {
				require.Equal(t, fix.groupID, got.ID())
				require.True(t, got.HasMember(fix.userID))
				return nil
			})

		// WHEN we add a user to the group
		err := fix.app.AddUserToGroup(context.Background(), fix.userID, fix.groupID)

		// THEN we get success and the user has been added to the group (see the GIVEN-THEN above)
		require.NoError(t, err)
	})

	t.Run("groupRepo load error", func(t *testing.T) {
		t.Parallel()

		fix := struct {
			*fixture
		}{
			fixture: newFixture(t),
		}

		// GIVEN a groupRepo that fails to get
		cause := errors.New("some_repo_error")
		fix.store.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(nil, cause)

		// WHEN we add a user to the group
		err := fix.app.AddUserToGroup(context.Background(), "irrelevant_user_id", "irrelevant_group_id")

		// THEN we get the error we expect
		require.Error(t, err)
		require.ErrorContains(t, err, cause.Error())
	})

	t.Run("groupRepo update error", func(t *testing.T) {
		t.Parallel()

		fix := struct {
			*fixture
		}{
			fixture: newFixture(t),
		}

		// GIVEN a groupRepo that loads an irrelevant group
		group := domain.NewGroup("irrelevant_group_id", "irrelevant_owner_id")
		fix.store.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(group, nil)

		// GIVEN a groupRepo that fails to update the group
		cause := errors.New("some_store_error")
		fix.store.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(cause)

		// WHEN we add a user to the group
		err := fix.app.AddUserToGroup(context.Background(), "irrelevant_user_id", "irrelevant_group_id")

		// THEN we get the error we expect
		require.Error(t, err)
		require.ErrorContains(t, err, cause.Error())
	})

	t.Run("full group", func(t *testing.T) {
		t.Parallel()

		// memberID create test member ids: member_id_0, member_id_1...
		memberID := func(n int) string { return fmt.Sprintf("member_id_%d", n) }

		fix := struct {
			*fixture
			groupID string
			ownerID string
		}{
			fixture: newFixture(t),
			groupID: "group_id",
			ownerID: "owner_id",
		}

		// GIVEN a groupRepo that loads a full group
		var fullGroup *domain.Group
		{
			fullGroup = domain.NewGroup(fix.groupID, fix.ownerID)
			for i := range domain.MaxMembers - 1 {
				err := fullGroup.AddMember(memberID(i))
				require.NoErrorf(t, err, "adding member %d", i)
			}
		}
		fix.store.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(fullGroup, nil)

		// WHEN we add a user to the group
		err := fix.app.AddUserToGroup(context.Background(), "new_user_id", fix.groupID)

		// THEN we get the error ErrGroupFull
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrGroupFull)
	})
}
