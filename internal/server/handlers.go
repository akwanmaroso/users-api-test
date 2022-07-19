package server

import (
	"net/http"

	authHttp "github.com/akwanmaroso/users-api/internal/auth/delivery/http"
	authUseCase "github.com/akwanmaroso/users-api/internal/auth/usecase"
	userHttp "github.com/akwanmaroso/users-api/internal/user/delivery/http"
	userRepository "github.com/akwanmaroso/users-api/internal/user/repository"
	userUseCase "github.com/akwanmaroso/users-api/internal/user/usecase"

	apiMiddlewares "github.com/akwanmaroso/users-api/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (s *Server) MapHandlers(e *echo.Echo) error {

	userDB := s.db.Database("users")
	userColl := userDB.Collection("users")

	userRepo := userRepository.NewMongoRepository(userColl)

	userUC := userUseCase.NewUserUseCase(s.cfg, userRepo, s.logger)
	authUC := authUseCase.NewAuthUseCase(s.cfg, userRepo, s.logger)

	userHandlers := userHttp.NewUserHandlers(s.cfg, userUC, s.logger)
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(s.cfg, userRepo, []string{"*"}, s.logger)

	v1 := e.Group("/api/v1")
	health := v1.Group("/health")
	userGroup := v1.Group("/users")
	authGroup := v1.Group("/auth")

	userHttp.MapUserRoutes(userGroup, userHandlers, mw)
	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)

	health.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OKE"})
	})

	return nil
}
