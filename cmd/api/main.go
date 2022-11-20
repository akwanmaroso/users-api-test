package main

import (
	"log"
	"os"
	"time"

	"github.com/akwanmaroso/users-api/config"
	"github.com/akwanmaroso/users-api/internal/server"
	"github.com/akwanmaroso/users-api/pkg/db"
	"github.com/akwanmaroso/users-api/pkg/logger"
	"github.com/akwanmaroso/users-api/pkg/utils"
	"golang.org/x/net/context"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Starting api server")
	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Load config: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewAPILogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	mongoDB, err := db.NewMongoDBConnection(ctx, cfg)
	if err != nil {
		appLogger.Fatalf("MongoDB init: %s", err)
	} else {
		appLogger.Infof("MongoDB connected, status: %#v")
	}

	redisClient, err := db.NewRedisClient(ctx, cfg)
	if err != nil {
		appLogger.Fatalf("Redis init: %s", err)
	}
	defer redisClient.Close()

	s := server.NewServer(cfg, mongoDB, redisClient, appLogger)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

}
