package usecase

import (
	"context"
	"fmt"

	"github.com/akwanmaroso/users-api/config"
	"github.com/akwanmaroso/users-api/internal/auth"
	"github.com/akwanmaroso/users-api/internal/models"
	"github.com/akwanmaroso/users-api/internal/user"
	"github.com/akwanmaroso/users-api/pkg/logger"
	"github.com/akwanmaroso/users-api/pkg/utils"
)

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

// Auth UseCase
type authUC struct {
	cfg      *config.Config
	userRepo user.Repository
	// redisRepo auth.RedisRepository
	// awsRepo   auth.AWSRepository
	logger logger.Logger
}

// Auth UseCase constructor
func NewAuthUseCase(cfg *config.Config, userRepo user.Repository, log logger.Logger) auth.UseCase {
	return &authUC{cfg: cfg, userRepo: userRepo, logger: log}
}

// Login user, returns user model with jwt token
func (u *authUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {

	foundUser, err := u.userRepo.GetByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, err
	}

	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		return nil, err
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

func (u *authUC) GenerateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, userID)
}
