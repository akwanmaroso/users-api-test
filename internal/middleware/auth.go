package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akwanmaroso/users-api/internal/models"
	"github.com/akwanmaroso/users-api/pkg/utils"
	"github.com/labstack/echo/v4"
)

func (mw *MiddlewareManager) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := utils.ExtractJWTFromRequest(c, mw.cfg.JWT.AccessSecret)
		if err != nil {
			mw.logger.Errorf("ExtractJWTFromRequest RequestID: %s,  Error: %s", utils.GetRequestID(c), err.Error())
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
		}

		fmt.Println(claims.ID)
		// id,claims.ID
		// fmt.Println(id)
		user, err := mw.userRepo.GetByID(c.Request().Context(), claims.ID)
		if err != nil {
			mw.logger.Errorf("GetByID RequestID: %s, Error: %s", utils.GetRequestID(c), err.Error())
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
		}

		c.Set("user", user)
		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, user)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func (mw *MiddlewareManager) RoleBasedAuthMiddleware(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*models.User)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": "Unauthorized",
				})
			}

			for _, role := range roles {
				if role == user.Role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    http.StatusForbidden,
				"message": "You don't have access to this endpoint",
			})
		}
	}
}
