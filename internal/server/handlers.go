package server

import (
	"net/http"

	"github.com/entl/evolyte-energy-provider-adapter/internal/config"
	"github.com/entl/evolyte-energy-provider-adapter/internal/enode"
	"github.com/entl/evolyte-energy-provider-adapter/internal/inverters"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func MapHandlers(s *echoServer) error {
	v1 := s.echoApp.Group("/api/v1")
	initalizeHealth(v1)
	initializeInverters(s.conf, v1, s.redisClient)

	return nil
}

func initalizeHealth(parentGroup *echo.Group) {
	parentGroup.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
}

func initializeInverters(conf *config.Config, parentGroup *echo.Group, redisClient *redis.Client) {
	authClient := enode.NewEnodeAuthClient(
		conf.Enode.ClientID,
		conf.Enode.ClientSecret,
		conf.Enode.OAuthBaseURL,
		conf.Enode.ApiURL,
		redisClient,
	)
	inverterClient := inverters.NewEnodeSolarInverterClient(
		authClient,
		conf.Enode.ApiURL,
		&http.Client{},
	)
	inverterUseCase := inverters.NewInverterUseCase(inverterClient, authClient)

	inverterHandler := inverters.NewInverterHandler(inverterUseCase)

	invertersGroup := parentGroup.Group("/enode/inverters")
	invertersGroup.GET("", inverterHandler.ListInverters)
}
