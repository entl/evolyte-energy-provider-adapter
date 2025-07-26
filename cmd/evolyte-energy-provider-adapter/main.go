package main

import (
	"log/slog"
	"os"

	"github.com/entl/evolyte-energy-provider-adapter/internal/config"
	"github.com/entl/evolyte-energy-provider-adapter/internal/server"
)

func main() {
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

	server.NewEchoServer(cfg).Start()
}
