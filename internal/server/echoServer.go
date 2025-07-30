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
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

type echoServer struct {
	echoApp     *echo.Echo
	conf        *config.Config
	redisClient *redis.Client
}

func NewEchoServer(conf *config.Config, redisClient *redis.Client) Server {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(echoLog.DEBUG)

	return &echoServer{
		echoApp:     echoApp,
		redisClient: redisClient,
		conf:        conf,
	}
}

func (s *echoServer) Start() error {
	s.echoApp.Use(middleware.Recover())
	s.echoApp.Use(middleware.Logger())
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", s.conf.Redis.Host, s.conf.Redis.Port),
		Password: s.conf.Redis.Password,
		DB:       s.conf.Redis.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
	}

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
