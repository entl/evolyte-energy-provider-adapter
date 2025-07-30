package inverters

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type InverterHandler struct {
	inverterUseCase *InverterUseCase
}

func NewInverterHandler(inverterUseCase *InverterUseCase) *InverterHandler {
	return &InverterHandler{
		inverterUseCase: inverterUseCase,
	}
}

func (h *InverterHandler) ListInverters(c echo.Context) error {
	after := c.QueryParam("after")
	before := c.QueryParam("before")
	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil {
		pageSize = 0 // Default to 0 if parsing fails
	}

	inverters, err := h.inverterUseCase.ListInverters(after, before, pageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list inverters")
	}

	return c.JSON(http.StatusOK, inverters)
}

func (h *InverterHandler) ListUserInverters(c echo.Context) error {
	userID := c.Param("userID")
	after := c.QueryParam("after")
	before := c.QueryParam("before")
	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil {
		pageSize = 0 // Default to 0 if parsing fails
	}

	inverters, err := h.inverterUseCase.ListUserInverters(userID, after, before, pageSize)
	if err != nil {
		slog.Error("Failed to list user inverters", "userID", userID, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list user inverters")
	}

	return c.JSON(http.StatusOK, inverters)
}

func (h *InverterHandler) GetInverter(c echo.Context) error {
	inverterID := c.Param("inverterID")
	inverter, err := h.inverterUseCase.GetInverter(inverterID)
	if err != nil {
		slog.Error("Failed to get inverter", "inverterID", inverterID, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get inverter")
	}

	return c.JSON(http.StatusOK, inverter)
}

func (h *InverterHandler) GetInverterProductionStatistics(c echo.Context) error {
	inverterID := c.Param("inverterID")
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		slog.Error("Invalid year parameter", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid year parameter")
	}
	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		slog.Error("Invalid month parameter", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid month parameter")
	}
	day, err := strconv.Atoi(c.Param("day"))
	if err != nil {
		slog.Error("Invalid day parameter", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid day parameter")
	}

	stats, err := h.inverterUseCase.GetInverterProductionStatistics(inverterID, year, month, day)
	if err != nil {
		slog.Error("Failed to get inverter production statistics", "inverterID", inverterID, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get inverter production statistics")
	}

	return c.JSON(http.StatusOK, stats)
}
