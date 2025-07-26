package server

import (
	"github.com/entl/evolyte-energy-provider-adapter/internal/config"
	"github.com/entl/evolyte-energy-provider-adapter/internal/enode"
	"github.com/labstack/echo/v4"
)

func MapHandlers(s *echoServer) error {
	v1 := s.echoApp.Group("/api/v1")
	initalizeHealth(v1)
	initializeEnodeAuth(s.conf, v1)

	return nil
}

func initalizeHealth(parentGroup *echo.Group) {
	parentGroup.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
}

func initializeEnodeAuth(conf *config.Config, parentGroup *echo.Group) {
	authClient := enode.NewEnodeAuthClient(
		conf.Enode.ClientID,
		conf.Enode.ClientSecret,
		conf.Enode.OAuthBaseURL,
		conf.Enode.ApiURL,
	)

	authHandler := enode.NewEnodeAuthHandler(authClient)

	authGroup := parentGroup.Group("/enode/auth")
	authGroup.GET("/authenticate", authHandler.Authenticate)
}
