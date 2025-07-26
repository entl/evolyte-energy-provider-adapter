package enode

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type enodeAuthHandler struct {
	AuthClient *enodeAuthClient
}

func NewEnodeAuthHandler(authClient *enodeAuthClient) *enodeAuthHandler {
	return &enodeAuthHandler{
		AuthClient: authClient,
	}
}

func (h *enodeAuthHandler) Authenticate(c echo.Context) error {
	res, err := h.AuthClient.Authenticate(c)
	if err != nil {
		slog.Error("Failed to authenticate with Enode", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to authenticate with Enode")
	}
	slog.Info("Successfully authenticated with Enode", "access_token", res.AccessToken, "expires_in", res.ExpiresIn)
	return c.JSON(200, res)
}
