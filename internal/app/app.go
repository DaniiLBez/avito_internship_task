package app

import (
	"fmt"
	"github.com/DaniiLBez/avito_internship_task/config"
	v1 "github.com/DaniiLBez/avito_internship_task/internal/controller/http/v1"
	"github.com/DaniiLBez/avito_internship_task/internal/repo"
	"github.com/DaniiLBez/avito_internship_task/internal/service"
	"github.com/DaniiLBez/avito_internship_task/pkg/hasher"
	"github.com/DaniiLBez/avito_internship_task/pkg/httpserver"
	"github.com/DaniiLBez/avito_internship_task/pkg/postgres"
	"github.com/DaniiLBez/avito_internship_task/pkg/validator"
	"github.com/labstack/echo/v4"
	"os"
	"os/signal"
	"syscall"
)

func Run(configPath string) {
	// Logger
	log := SetLogger(InfoLevel)

	// Configuration
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Error("Config error: %s", err)
	}

	// Repositories
	log.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Error(fmt.Sprintf("app - Run - pgdb.NewServices: %w", err))
	}
	defer pg.Close()

	// Repositories
	log.Info("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	// Services dependencies
	log.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos:    repositories,
		Hasher:   hasher.NewSHA1Hasher(cfg.Hasher.Salt),
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}
	services := service.NewServices(&deps)

	// Echo handler
	log.Info("Initializing handlers and routes...")
	handler := echo.New()
	// setup handler validator as lib validator
	handler.Validator = validator.NewCustomValidator()
	v1.NewRouter(handler, services)

	// HTTP server
	log.Info("Starting http server...")
	log.Debug(fmt.Sprintf("Server port: %s", cfg.HTTP.Port))
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Sprintf("app - Run - httpServer.Notify: %w", err))
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %w", err))
	}
}
