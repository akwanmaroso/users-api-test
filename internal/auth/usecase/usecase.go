package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/akwanmaroso/users-api/config"
	"github.com/akwanmaroso/users-api/internal/auth"
	"github.com/akwanmaroso/users-api/internal/models"
	"github.com/akwanmaroso/users-api/internal/user"
	"github.com/akwanmaroso/users-api/pkg/logger"
	"github.com/akwanmaroso/users-api/pkg/utils"
)

const (
	accessTokenDuration  = time.Duration(time.Hour * 1)
	refreshTokenDuration = time.Duration(time.Hour * 24)
)

// Auth UseCase
type authUC struct {
	cfg      *config.Config
	userRepo user.Repository
	// redisRepo auth.RedisRepository
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

	accessToken, err := utils.GenerateJWTToken(foundUser, u.cfg.AccessSecret, accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateJWTToken(foundUser, u.cfg.RefreshSecret, refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &models.UserWithToken{
		User:         foundUser,
		AccesssToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *authUC) Refresh(ctx context.Context, token string) (*models.UserWithToken, error) {
	claims, err := utils.ExtractUser(token, u.cfg.RefreshSecret)
	if err != nil {
		return nil, err
	}

	if claims != nil {
		foundUser, err := u.userRepo.GetByID(ctx, claims.ID)
		if err != nil {
			return nil, err
		}

		accessToken, err := utils.GenerateJWTToken(foundUser, u.cfg.AccessSecret, accessTokenDuration)
		if err != nil {
			return nil, err
		}

		refreshToken, err := utils.GenerateJWTToken(foundUser, u.cfg.RefreshSecret, accessTokenDuration)
		if err != nil {
			return nil, err
		}

		return &models.UserWithToken{
			User:         foundUser,
			AccesssToken: accessToken,
			RefreshToken: refreshToken,
		}, nil
	} else {
		return nil, errors.New("Invalid token")
	}
}

// func (u *authUC) GenerateUserKey(userID string) string {
// 	return fmt.Sprintf("%s: %s", basePrefix, userID)
// }
