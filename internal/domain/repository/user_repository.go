

package repository

import (
    "context"
	"building-report-backend/internal/domain/entity"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id string) error
    FindByID(ctx context.Context, id string) (*entity.User, error)
    FindByUsername(ctx context.Context, username string) (*entity.User, error)
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
    FindAll(ctx context.Context, limit, offset int) ([]*entity.User, int64, error)
}