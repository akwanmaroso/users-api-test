package http

import (
	"net/http"

	"github.com/akwanmaroso/users-api/config"
	"github.com/akwanmaroso/users-api/internal/models"
	"github.com/akwanmaroso/users-api/internal/user"
	"github.com/akwanmaroso/users-api/pkg/logger"
	"github.com/akwanmaroso/users-api/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userHandlersImpl struct {
	cfg    *config.Config
	userUC user.UseCase
	logger logger.Logger
}

func NewUserHandlers(cfg *config.Config, userUC user.UseCase, logger logger.Logger) user.Handlers {
	return &userHandlersImpl{cfg: cfg, userUC: userUC, logger: logger}
}

func (h *userHandlersImpl) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		n := &models.User{}
		if err := c.Bind(n); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":  http.StatusBadRequest,
				"error": err.Error(),
			})
		}

		createdUser, err := h.userUC.Create(ctx, n)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":  http.StatusInternalServerError,
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"code": http.StatusOK,
			"data": createdUser,
		})
	}
}

func (h *userHandlersImpl) Current() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*models.User)
		if !ok {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
		}

		user.SanitizePassword()

		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": http.StatusOK,
			"data": user,
		})
	}
}

func (h *userHandlersImpl) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		list, err := h.userUC.List(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":  http.StatusInternalServerError,
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": http.StatusOK,
			"data": list,
		})
	}
}

func (h *userHandlersImpl) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		id := c.Param("id")

		n := &models.User{}
		err := c.Bind(n)
		if err != nil {
			h.logger.Error(err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":  http.StatusBadRequest,
				"error": err.Error(),
			})
		}

		n.ID, err = primitive.ObjectIDFromHex(id)
		updatedNews, err := h.userUC.Update(ctx, n)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":  http.StatusInternalServerError,
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, updatedNews)
	}
}

func (h *userHandlersImpl) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		id := c.Param("id")

		err := h.userUC.Delete(ctx, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":  http.StatusInternalServerError,
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusOK,
			"message": "Success deleted user",
		})
	}
}
