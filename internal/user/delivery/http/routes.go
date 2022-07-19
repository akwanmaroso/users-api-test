package http

import (
	"github.com/akwanmaroso/users-api/internal/middleware"
	"github.com/akwanmaroso/users-api/internal/user"
	"github.com/labstack/echo/v4"
)

func MapUserRoutes(userGroup *echo.Group, h user.Handlers, mw *middleware.MiddlewareManager) {
	userGroup.POST("", h.Create())
	userGroup.DELETE("/:id", h.Delete(), mw.AuthMiddleware, mw.RoleBasedAuthMiddleware([]string{"admin"}))
	userGroup.GET("/current", h.Detail(), mw.AuthMiddleware)
	userGroup.GET("", h.List(), mw.AuthMiddleware, mw.RoleBasedAuthMiddleware([]string{"admin"}))
}
