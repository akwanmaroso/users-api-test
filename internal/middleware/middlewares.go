package middleware

import (
	"github.com/akwanmaroso/users-api/config"
	"github.com/akwanmaroso/users-api/internal/user"
	"github.com/akwanmaroso/users-api/pkg/logger"
)

type MiddlewareManager struct {
	cfg      *config.Config
	userRepo user.Repository
	origins  []string
	logger   logger.Logger
}

func NewMiddlewareManager(cfg *config.Config, userRepo user.Repository, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{
		cfg:      cfg,
		origins:  origins,
		userRepo: userRepo,
		logger:   logger,
	}
}
