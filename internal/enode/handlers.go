package enode

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type enodeAuthHandler struct {
	AuthClient *EnodeAuthClient
}

func NewEnodeAuthHandler(authClient *EnodeAuthClient) *enodeAuthHandler {
	return &enodeAuthHandler{
		AuthClient: authClient,
	}
}

func (h *enodeAuthHandler) Authenticate(c echo.Context) error {
	res, err := h.AuthClient.authenticate()
	if err != nil {
		slog.Error("Failed to authenticate with Enode", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to authenticate with Enode")
	}
	slog.Debug("Successfully authenticated with Enode")
	return c.JSON(200, res)
}
