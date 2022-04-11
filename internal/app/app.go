// Package app configures and runs application.
package app

import (
	"finstar-test-task/config"
	v1 "finstar-test-task/internal/controller/http/v1"
	"finstar-test-task/internal/usecase/repo"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"finstar-test-task/pkg/httpserver"
	"finstar-test-task/pkg/logger"
	"finstar-test-task/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(GetDSN(cfg), cfg.PG.LogLevel, cfg.PG.AutoMigrate, cfg.PG.GenerateSeeds)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	httpServer := httpserver.New(v1.New(repo.New(pg), l), httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()

	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

func GetDSN(c *config.Config) string {
	return fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s sslmode=%s",
		c.PG.Host,
		c.PG.Port,
		c.PG.User,
		c.PG.Name,
		c.PG.Password,
		c.PG.SSLMode,
	)
}
