package usecase

import (
	"context"

	"github.com/akwanmaroso/users-api/config"
	"github.com/akwanmaroso/users-api/internal/models"
	"github.com/akwanmaroso/users-api/internal/user"
	"github.com/akwanmaroso/users-api/pkg/logger"
	"github.com/akwanmaroso/users-api/pkg/utils"
)

type userUseCaseImpl struct {
	cfg      *config.Config
	userRepo user.Repository
	logger   logger.Logger
}

func NewUserUseCase(cfg *config.Config, userRepo user.Repository, logger logger.Logger) user.UseCase {
	return &userUseCaseImpl{
		cfg:      cfg,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (u *userUseCaseImpl) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if err := user.PrepareCreate(); err != nil {
		return nil, err
	}

	result, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	result.SanitizePassword()

	return result, nil
}

func (u *userUseCaseImpl) Delete(ctx context.Context, id string) error {

	err := u.userRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCaseImpl) Update(ctx context.Context, user *models.User) (*models.User, error) {
	user.PrepareUpdate()

	if err := utils.ValidateIsOwner(ctx, user.ID.Hex(), u.logger); err != nil {
		return nil, err
	}

	result, err := u.userRepo.Edit(ctx, user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *userUseCaseImpl) List(ctx context.Context) ([]*models.User, error) {
	result, err := u.userRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(result); i++ {
		result[i].SanitizePassword()
	}

	return result, nil
}
