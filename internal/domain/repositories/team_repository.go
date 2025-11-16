package repositories

import (
	"context"
)

type TeamRepository interface {
	Create(ctx context.Context, name string) error
	GetByName(ctx context.Context, name string) (string, error)
	List(ctx context.Context) ([]string, error)
}
