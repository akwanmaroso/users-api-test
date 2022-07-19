package utils

import (
	"context"
	"errors"

	"github.com/akwanmaroso/users-api/internal/models"
	"github.com/labstack/echo/v4"
)

type ReqIDCtxKey struct{}

func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func GetRequestCtx(c echo.Context) context.Context {
	return context.WithValue(c.Request().Context(), ReqIDCtxKey{}, GetRequestID(c))
}

type UserCtxKey struct{}

// Get user from context
func GetUserFromCtx(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value(UserCtxKey{}).(*models.User)
	if !ok {
		return nil, errors.New("Unauthorized")
	}

	return user, nil
}

// Read request body and validate
func ReadRequest(ctx echo.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(ctx.Request().Context(), request)
}
