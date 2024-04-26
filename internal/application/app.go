package application

import (
	"context"
	"fmt"
	"time"

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

	if err := a.store.Create(ctx, group); err != nil {
		return "", fmt.Errorf("creating: %v", err)
	}

	return groupID, nil

}

func (a *App) GetGroup(ctx context.Context, groupID string) (*domain.Group, error) {
	var group *domain.Group

	group, err := a.store.Load(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("creating: %v", err)
	}

	return group, nil
}

func (a *App) AddUserToGroup(ctx context.Context, userID, groupID string, options ...Option) error {
	do := func(ctx context.Context) error {
		group, err := a.store.Load(ctx, groupID)
		if err != nil {
			return fmt.Errorf("creating: %v", err)
		}

		if err := group.AddMember(userID); err != nil {
			return fmt.Errorf("adding: %w", err)
		}

		if d, ok := delayBeforeUpdating(options...); ok {
			time.Sleep(d)
		}

		if err := a.store.Update(ctx, group); err != nil {
			return fmt.Errorf("updating: %v", err)
		}

		return nil
	}

	return do(ctx)
}

func delayBeforeUpdating(options ...Option) (time.Duration, bool) {
	for _, o := range options {
		if raw, ok := o.(DelayBeforeUpdating); ok {
			return time.Duration(raw), true
		}
	}

	return time.Duration(0), false
}
