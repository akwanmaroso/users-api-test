package user

import (
	"context"

	"github.com/akwanmaroso/users-api/internal/models"
)

type Repository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	Edit(ctx context.Context, user *models.User) (*models.User, error)
	List(ctx context.Context) ([]*models.User, error)
	Delete(ctx context.Context, id string) error
}
