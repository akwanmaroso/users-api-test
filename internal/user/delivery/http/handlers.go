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

// Create godoc
// @Summary Create users
// @Description Create users handler
// @Tags Users
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Router /users [post]
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

// Create godoc
// @Summary Create users
// @Description Create users handler
// @Tags Users
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Router /users [post]
func (h *userHandlersImpl) Detail() echo.HandlerFunc {
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

// Create godoc
// @Summary Create users
// @Description Create users handler
// @Tags Users
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Router /users [post]
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

// Update godoc
// @Summary Update news
// @Description Update news handler
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "news_id"
// @Success 200 {object} models.News
// @Router /news/{id} [put]
func (h *userHandlersImpl) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := utils.GetRequestCtx(c)
		id := c.Param("id")

		n := &models.User{}
		err := c.Bind(n)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":  http.StatusBadRequest,
				"error": err.Error(),
			})
		}

		n.ID, err = primitive.ObjectIDFromHex(id)
		if err := c.Bind(n); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":  http.StatusBadRequest,
				"error": err.Error(),
			})
		}

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

// Delete godoc
// @Summary Delete users
// @Description Delete by id users handler
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {string} string	"ok"
// @Router /users/{id} [delete]
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
