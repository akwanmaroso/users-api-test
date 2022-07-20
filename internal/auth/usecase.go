package auth

import (
	"context"

	"github.com/akwanmaroso/users-api/internal/models"
)

// Auth repository interface
type UseCase interface {
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Refresh(ctx context.Context, token string) (*models.UserWithToken, error)
}
