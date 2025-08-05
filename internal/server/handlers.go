package server

import (
	"net/http"

	"github.com/entl/evolyte-energy-provider-adapter/internal/enode"
	"github.com/entl/evolyte-energy-provider-adapter/internal/inverters"
	"github.com/labstack/echo/v4"
)

func MapHandlers(s *echoServer) error {
	v1 := s.echoApp.Group("/api/v1")
	initalizeHealth(v1)
	initializeInverters(s, v1)

	return nil
}

func initalizeHealth(parentGroup *echo.Group) {
	parentGroup.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
}

func initializeInverters(s *echoServer, parentGroup *echo.Group) {
	authClient := enode.NewEnodeAuthClient(
		s.conf.Enode.ClientID,
		s.conf.Enode.ClientSecret,
		s.conf.Enode.OAuthBaseURL,
		s.conf.Enode.ApiURL,
		s.redisClient,
	)
	inverterClient := inverters.NewEnodeSolarInverterClient(
		authClient,
		s.conf.Enode.ApiURL,
		&http.Client{},
		s.inverterQueries,
	)
	inverterUseCase := inverters.NewInverterUseCase(inverterClient, authClient, s.inverterQueries, s.validator)

	inverterHandler := inverters.NewInverterHandler(inverterUseCase)

	invertersGroup := parentGroup.Group("/enode/inverters")
	userInvertersGroup := parentGroup.Group("/enode/users")

	userInvertersGroup.GET("/:userID", inverterHandler.ListUserInverters)
	userInvertersGroup.POST("/:userID/link", inverterHandler.LinkInverter)

	invertersGroup.GET("", inverterHandler.ListInverters)
	invertersGroup.GET("/:inverterID", inverterHandler.GetInverter)
	invertersGroup.GET("/:inverterID/stats", inverterHandler.GetInverterProductionStatistics)
	invertersGroup.POST("", inverterHandler.AddInverter)
}
