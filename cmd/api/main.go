package main

import (
	"fmt"
	"log"

	"github.com/firdanbash/go-clean-boiler/internal/domain"
	"github.com/firdanbash/go-clean-boiler/internal/handler"
	"github.com/firdanbash/go-clean-boiler/internal/repository/postgres"
	"github.com/firdanbash/go-clean-boiler/internal/router"
	"github.com/firdanbash/go-clean-boiler/internal/service"
	"github.com/firdanbash/go-clean-boiler/pkg/config"
	"github.com/firdanbash/go-clean-boiler/pkg/database"
	"github.com/firdanbash/go-clean-boiler/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.Log.Level, cfg.Log.Encoding); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting application",
		zap.String("app", cfg.App.Name),
		zap.String("env", cfg.App.Env),
	)

	// Initialize database
	if err := database.Init(cfg); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	// Auto migrate models
	if err := database.AutoMigrate(&domain.User{}); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}
	logger.Info("Database migrations completed successfully")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(database.DB)

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.Expiration.String())

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	r := router.SetupRouter(authHandler, userHandler, cfg.JWT.Secret)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	logger.Info("Server starting", zap.String("address", addr))

	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
