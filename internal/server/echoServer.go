package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/entl/evolyte-energy-provider-adapter/internal/config"
	"github.com/entl/evolyte-energy-provider-adapter/internal/db"
	"github.com/entl/evolyte-energy-provider-adapter/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

type echoServer struct {
	echoApp         *echo.Echo
	redisClient     *redis.Client
	inverterQueries *db.Queries
	conf            *config.Config
	validator       *utils.CustomValidator
}

func NewEchoServer(conf *config.Config, redisClient *redis.Client, inverterQueries *db.Queries, validator *utils.CustomValidator) Server {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(echoLog.DEBUG)

	return &echoServer{
		echoApp:         echoApp,
		redisClient:     redisClient,
		inverterQueries: inverterQueries,
		conf:            conf,
		validator:       validator,
	}
}

func (s *echoServer) Start() error {
	s.echoApp.Use(middleware.Recover())
	s.echoApp.Use(middleware.Logger())
	s.echoApp.Validator = s.validator

	go func() {
		slog.Info("Starting Echo server", "port", s.conf.Server.Port)
		if err := s.echoApp.Start(fmt.Sprintf(":%s", s.conf.Server.Port)); err != nil {
			slog.Error("Failed to start server", "error", err)
		}
	}()

	if err := MapHandlers(s); err != nil {
		slog.Error("Failed to map handlers", "error", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	slog.Info("Shutting down server gracefully")
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)

	defer shutdown()
	return s.echoApp.Shutdown(ctx)
}
