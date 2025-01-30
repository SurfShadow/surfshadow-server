package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/SurfShadow/surfshadow-server/docs"
	"github.com/SurfShadow/surfshadow-server/internal/infrastructure/config"
	"github.com/SurfShadow/surfshadow-server/internal/infrastructure/db"
	"github.com/SurfShadow/surfshadow-server/internal/infrastructure/migrations"
	"github.com/SurfShadow/surfshadow-server/internal/presentation/server"
	"github.com/SurfShadow/surfshadow-server/pkg/logger"
)

// @title SurfShadow API
// @version 1.0.0
// @description Rent VPN VLESS configs
// @BasePath /api/v1
// @servers.url https://api.surfshadow.com/api/v1
// @servers.description Production server
// @servers.url /api/v1
// @servers.description Development server.
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	err = logger.NewLogger(cfg.Logger)
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}

	logger.Instance.Debug("Logger initialized successfully")

	logger.Instance.Infof("Loaded config of server: %+v", cfg.Server)
	logger.Instance.Infof("Loaded config of database: %+v", cfg.DB)
	logger.Instance.Infof("Loaded config of logger: %+v", cfg.Logger)
	logger.Instance.Infof("Loaded config of metrics: %+v", cfg.Metrics)
	logger.Instance.Infof("Loaded config of auth API: %+v", cfg.AuthAPI)

	logger.Instance.Info("Database connecting...")

	database, err := db.NewPsqlDB(cfg.DB)
	if err != nil {
		logger.Instance.Fatalf("Failed to connect to database: %v", err)
	}

	defer func() {
		err = database.Close()
		if err != nil {
			logger.Instance.Errorf("Failed to close database connection: %v", err)
		}
	}()

	logger.Instance.Info("Database connected successfully")

	migrationsPath := "file://migrations"
	logger.Instance.Infof("Applying migrations from %s", migrationsPath)

	err = migrations.ApplyMigrations(database, migrationsPath)
	if err != nil {
		logger.Instance.Fatalf("Failed to apply migrations: %v", err)
	}

	logger.Instance.Info("Migrations applied successfully")

	serv := server.NewServer(cfg.Server, database)
	go func() {
		if err = serv.Start(); err != nil {
			logger.Instance.Fatalf("Failed to start server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Instance.Info("Shutting down server...")

	if err = serv.Stop(); err != nil {
		logger.Instance.Errorf("Error while stopping server: %v", err)
	} else {
		logger.Instance.Info("Server stopped gracefully")
	}
}
