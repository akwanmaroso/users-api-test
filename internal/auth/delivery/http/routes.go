package http

import (
	"github.com/akwanmaroso/users-api/internal/auth"
	"github.com/akwanmaroso/users-api/internal/middleware"
	"github.com/labstack/echo/v4"
)

// Map auth routes
func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mw *middleware.MiddlewareManager) {
	authGroup.POST("/login", h.Login())
	authGroup.POST("/refresh", h.Refresh())
}
