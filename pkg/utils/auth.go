package utils

import (
	"context"
	"errors"

	"github.com/akwanmaroso/users-api/pkg/logger"
)

func ValidateIsOwner(ctx context.Context, creatorID string, logger logger.Logger) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return err
	}

	if user.ID.Hex() != creatorID && user.Role != "admin" {
		return errors.New("Unauthorized")
	}

	return nil
}
