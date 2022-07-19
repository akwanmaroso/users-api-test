package user

import (
	"context"

	"github.com/akwanmaroso/users-api/internal/models"
)

type UseCase interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user *models.User) (*models.User, error)
	List(ctx context.Context) ([]*models.User, error)
}
