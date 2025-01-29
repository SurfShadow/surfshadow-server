package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/SurfShadow/surfshadow-server/internal/application/usecases"
	"github.com/SurfShadow/surfshadow-server/internal/infrastructure/config"
	"github.com/SurfShadow/surfshadow-server/internal/infrastructure/db/repositories"
	"github.com/SurfShadow/surfshadow-server/internal/presentation/handlers"
	"github.com/SurfShadow/surfshadow-server/internal/presentation/routes"
	"github.com/SurfShadow/surfshadow-server/pkg/logger"
)

const (
	privateKeyPath = "ssl/server.crt"
	publicKeyPath  = "ssl/server.pem"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.ServerConfig, db *sqlx.DB) *Server {
	logger.Instance.Debug("Initializing Server")

	proxyClientRepo := repositories.NewProxyClientRepository(db)
	proxyClientUseCase := usecases.NewProxyClientUseCase(proxyClientRepo)
	proxyClientHandler := handlers.NewProxyClientHandler(proxyClientUseCase)

	router := routes.InitRoutes(proxyClientHandler)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	logger.Instance.Infof("Server initialized with address: %s", httpServer.Addr)
	logger.Instance.Debugf("Server read timeout: %s", httpServer.ReadTimeout)
	logger.Instance.Debugf("Server write timeout: %s", httpServer.WriteTimeout)
	logger.Instance.Debugf("Server idle timeout: %s", httpServer.IdleTimeout)

	return &Server{
		httpServer: httpServer,
	}
}

func (s *Server) Start() error {
	logger.Instance.Infof("Starting HTTP server on %s", s.httpServer.Addr)
	logger.Instance.Debug("HTTP server starting...")

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Instance.Errorf("Failed to start HTTP server: %v", err)
		return err
	}

	logger.Instance.Info("HTTP server started successfully")

	return nil
}

func (s *Server) Stop() error {
	logger.Instance.Info("Stopping server")
	logger.Instance.Debug("HTTP server stopping...")

	if err := s.httpServer.Close(); err != nil {
		logger.Instance.Errorf("Failed to stop HTTP server: %v", err)
		return err
	}

	logger.Instance.Info("HTTP server stopped successfully")

	return nil
}
