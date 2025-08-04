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

	inverters, err := h.inverterUseCase.ListInverters(c.Request().Context(), after, before, pageSize)
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

	inverters, err := h.inverterUseCase.ListUserInverters(c.Request().Context(), userID, after, before, pageSize)
	if err != nil {
		slog.Error("Failed to list user inverters", "userID", userID, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list user inverters")
	}

	return c.JSON(http.StatusOK, inverters)
}

func (h *InverterHandler) GetInverter(c echo.Context) error {
	inverterID := c.Param("inverterID")
	inverter, err := h.inverterUseCase.GetInverter(c.Request().Context(), inverterID)
	if err != nil {
		slog.Error("Failed to get inverter", "inverterID", inverterID, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get inverter")
	}

	return c.JSON(http.StatusOK, inverter)
}

func (h *InverterHandler) GetInverterProductionStatistics(c echo.Context) error {
	inverterID := c.Param("inverterID")
	year, err := strconv.Atoi(c.QueryParam("year"))
	if err != nil {
		slog.Error("Invalid year parameter", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid year parameter")
	}
	month, err := strconv.Atoi(c.QueryParam("month"))
	if err != nil {
		slog.Error("Invalid month parameter", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid month parameter")
	}
	day, err := strconv.Atoi(c.QueryParam("day"))
	if err != nil {
		slog.Error("Invalid day parameter", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid day parameter")
	}

	stats, err := h.inverterUseCase.GetInverterProductionStatistics(c.Request().Context(), inverterID, year, month, day)
	if err != nil {
		slog.Error("Failed to get inverter production statistics", "inverterID", inverterID, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get inverter production statistics")
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *InverterHandler) AddInverter(c echo.Context) error {
	var request AddInverterRequest
	if err := c.Bind(&request); err != nil {
		slog.Error("Failed to bind request", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}
	if err := c.Validate(request); err != nil {
		slog.Error("Validation failed for AddInverterRequest", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed")
	}

	response, err := h.inverterUseCase.AddInverter(c.Request().Context(), request)
	if err != nil {
		slog.Error("Failed to add inverter", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add inverter")
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *InverterHandler) LinkInverter(c echo.Context) error {
	userID := c.Param("userID")
	var request LinkInverterRequest
	if err := c.Bind(&request); err != nil {
		slog.Error("Failed to bind request", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}
	if err := c.Validate(request); err != nil {
		slog.Error("Validation failed for LinkInverterRequest", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed")
	}

	response, err := h.inverterUseCase.LinkInverter(c.Request().Context(), userID, request)
	if err != nil {
		slog.Error("Failed to link inverter", "userID", userID, "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to link inverter")
	}

	return c.JSON(http.StatusOK, response)
}
