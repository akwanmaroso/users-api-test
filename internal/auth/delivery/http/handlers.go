package http

import (
	"net/http"

	"github.com/akwanmaroso/users-api/config"
	"github.com/akwanmaroso/users-api/internal/auth"
	"github.com/akwanmaroso/users-api/internal/models"
	"github.com/akwanmaroso/users-api/pkg/logger"
	"github.com/akwanmaroso/users-api/pkg/utils"
	"github.com/labstack/echo/v4"
)

type authHandlersImpl struct {
	cfg    *config.Config
	authUC auth.UseCase
	logger logger.Logger
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, logger logger.Logger) auth.Handlers {
	return &authHandlersImpl{cfg: cfg, authUC: authUC, logger: logger}
}

func (h *authHandlersImpl) Login() echo.HandlerFunc {
	type Login struct {
		Username string `json:"username" validate:"required,lte=60"`
		Password string `json:"password,omitempty" validate:"required"`
	}
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":  http.StatusOK,
				"error": err,
			})
		}

		userWithToken, err := h.authUC.Login(ctx, &models.User{
			Username: login.Username,
			Password: login.Password,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":  http.StatusInternalServerError,
				"error": err,
			})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"code": http.StatusOK,
			"data": userWithToken,
		})
	}
}

func (h *authHandlersImpl) Refresh() echo.HandlerFunc {
	type Refresh struct {
		Token string `json:"token" validate:"required"`
	}
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)

		refresh := &Refresh{}
		if err := utils.ReadRequest(c, refresh); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":  http.StatusOK,
				"error": err,
			})
		}

		userWithToken, err := h.authUC.Refresh(ctx, refresh.Token)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":  http.StatusInternalServerError,
				"error": err,
			})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"code": http.StatusOK,
			"data": userWithToken,
		})
	}
}
