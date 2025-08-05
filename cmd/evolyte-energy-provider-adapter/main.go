package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/entl/evolyte-energy-provider-adapter/internal/config"
	"github.com/entl/evolyte-energy-provider-adapter/internal/db"
	"github.com/entl/evolyte-energy-provider-adapter/internal/server"
	"github.com/entl/evolyte-energy-provider-adapter/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))

	slog.Info("Starting Evolyte Energy Provider Adapter")
	cfg, err := config.LoadConfig(".env.docker")
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		panic(err)
	}

	conn, err := pgx.Connect(ctx, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DB))
	if err != nil {
		slog.Error("Failed to connect to Postgres", "error", err)
		panic(err)
	}
	defer conn.Close(ctx)

	inverterQueries := db.New(conn)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if redisClient.Ping(ctx).Err() != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		panic(err)
	} else {
		slog.Info("Connected to Redis successfully")
	}
	defer redisClient.Close()

	slog.Info("Initializing Validator")
	val := utils.NewCustomValidator(validator.New())

	server.NewEchoServer(cfg, redisClient, inverterQueries, val).Start()
}
