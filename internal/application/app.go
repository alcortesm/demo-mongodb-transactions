package application

import (
	"context"
	"fmt"

	"github.com/alcortesm/demo-mongodb-transactions/internal/domain"
)

//go:generate mockgen -source=app.go -destination=mock_dependencies_test.go -package=application_test

type App struct {
	uuider Uuider
	store  Store
}

type Store interface {
	Create(ctx context.Context, group *domain.Group) error
	Update(ctx context.Context, group *domain.Group) error
	Load(ctx context.Context, id string) (*domain.Group, error)
}

// Uuider knows how to return V4 UUIDs.
type Uuider interface {
	NewString() string
}

func New(
	uuider Uuider,
	store Store,
) *App {
	return &App{
		uuider: uuider,
		store:  store,
	}
}

func (a *App) CreateGroup(ctx context.Context, ownerID string) (string, error) {
	groupID := a.uuider.NewString()
	group := domain.NewGroup(groupID, ownerID)

	do := func(ctx context.Context) error {
		if err := a.store.Create(ctx, group); err != nil {
			return fmt.Errorf("creating: %v", err)
		}

		return nil
	}

	if err := do(ctx); err != nil {
		return "", err
	}

	return groupID, nil

}

func (a *App) GetGroup(ctx context.Context, groupID string) (*domain.Group, error) {
	var group *domain.Group

	do := func(ctx context.Context) error {
		var err error

		group, err = a.store.Load(ctx, groupID)
		if err != nil {
			return fmt.Errorf("creating: %v", err)
		}

		return nil
	}

	if err := do(ctx); err != nil {
		return nil, err
	}

	return group, nil
}

func (a *App) AddUserToGroup(ctx context.Context, userID, groupID string) error {
	do := func(ctx context.Context) error {
		group, err := a.store.Load(ctx, groupID)
		if err != nil {
			return fmt.Errorf("creating: %v", err)
		}

		if err := group.AddMember(userID); err != nil {
			return fmt.Errorf("adding: %w", err)
		}

		if err := a.store.Update(ctx, group); err != nil {
			return fmt.Errorf("updating: %v", err)
		}

		return nil
	}

	return do(ctx)
}
