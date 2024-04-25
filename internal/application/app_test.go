package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alcortesm/demo-mongodb-transactions/internal/application"
	"github.com/alcortesm/demo-mongodb-transactions/internal/domain"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type fixture struct {
	app       *application.App
	uuider    *MockUuider
	groupRepo *MockGroupRepo
}

func newFixture(t *testing.T) *fixture {
	ctrl := gomock.NewController(t)
	uuider := NewMockUuider(ctrl)
	groupRepo := NewMockGroupRepo(ctrl)

	app := application.New(uuider, groupRepo)

	return &fixture{
		app:       app,
		uuider:    uuider,
		groupRepo: groupRepo,
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
			fix.groupRepo.EXPECT().
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
		fix.groupRepo.EXPECT().
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
		fix.groupRepo.EXPECT().
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
		fix.groupRepo.EXPECT().
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
		fix.groupRepo.EXPECT().
			Load(gomock.Any(), fix.groupID).
			Return(group, nil)

		// GIVEN-THEN a groupRepo expecting an Update with the new user as a member
		fix.groupRepo.EXPECT().
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
		fix.groupRepo.EXPECT().
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
		fix.groupRepo.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(group, nil)

		// GIVEN a groupRepo that fails to update the group
		cause := errors.New("some_store_error")
		fix.groupRepo.EXPECT().
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

		fix := struct {
			*fixture
		}{
			fixture: newFixture(t),
		}

		// GIVEN a groupRepo that loads a full group
		group := domain.NewGroup("irrelevant_group_id", "irrelevant_owner_id")
		t.Fatal("TODO return a full group from here")
		fix.groupRepo.EXPECT().
			Load(gomock.Any(), gomock.Any()).
			Return(group, nil)

		// GIVEN a groupRepo that fails to update the group
		cause := errors.New("some_store_error")
		fix.groupRepo.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(cause)

		// WHEN we add a user to the group
		err := fix.app.AddUserToGroup(context.Background(), "irrelevant_user_id", "irrelevant_group_id")

		// THEN we get the error we expect
		require.Error(t, err)
		require.ErrorContains(t, err, cause.Error())
	})
}
