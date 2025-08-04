package inverters

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/entl/evolyte-energy-provider-adapter/internal/db"
	"github.com/entl/evolyte-energy-provider-adapter/internal/enode"
	"github.com/entl/evolyte-energy-provider-adapter/internal/utils"
)

type InverterUseCase struct {
	inverterClient  SolarInverterClient
	authClient      *enode.EnodeAuthClient
	inverterQueries *db.Queries
	validator       *utils.CustomValidator
}

func NewInverterUseCase(inverterClient SolarInverterClient, authClient *enode.EnodeAuthClient, inverterQueries *db.Queries, validator *utils.CustomValidator) *InverterUseCase {
	return &InverterUseCase{
		inverterClient:  inverterClient,
		authClient:      authClient,
		inverterQueries: inverterQueries,
		validator:       validator,
	}
}

func (uc *InverterUseCase) ListInverters(ctx context.Context, after string, before string, pageSize int) (*SolarInverterResponse, error) {
	// Create context with timeout for token acquisition
	tokenCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	token, err := uc.authClient.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	bearerToken := "Bearer " + token

	inverters, err := uc.inverterClient.ListInverters(ctx, bearerToken, after, before, pageSize)
	// Check for timeout specifically
	if errors.Is(tokenCtx.Err(), context.DeadlineExceeded) {
		slog.Error("timeout while getting access token", "error", tokenCtx.Err())
		return nil, fmt.Errorf("timeout while getting access token: %w", tokenCtx.Err())
	}
	if err != nil {
		slog.Error("Failed to list inverters", "error", err)
		return nil, fmt.Errorf("failed to list inverters: %w", err)
	}
	return inverters, nil
}

func (uc *InverterUseCase) ListUserInverters(ctx context.Context, userID string, after string, before string, pageSize int) (*SolarInverterResponse, error) {
	token, err := uc.authClient.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	bearerToken := "Bearer " + token

	return uc.inverterClient.ListUserInverters(ctx, bearerToken, userID, after, before, pageSize)
}

func (uc *InverterUseCase) GetInverter(ctx context.Context, inverterID string) (*SolarInverter, error) {
	token, err := uc.authClient.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	bearerToken := "Bearer " + token

	inverter, err := uc.inverterClient.GetInverter(ctx, bearerToken, inverterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inverter: %w", err)
	}
	return inverter, nil
}

func (uc *InverterUseCase) GetInverterProductionStatistics(ctx context.Context, inverterID string, year int, month int, day int) (*InverterStatistic, error) {
	token, err := uc.authClient.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	bearerToken := "Bearer " + token

	params := InverterStatisticParams{
		Year:  year,
		Month: month,
		Day:   day,
	}
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid inverter statistic parameters: %w", err)
	}

	return uc.inverterClient.GetInverterProductionStatistics(ctx, bearerToken, inverterID, params)
}

func (uc *InverterUseCase) AddInverter(ctx context.Context, request AddInverterRequest) (*AddInverterResponse, error) {
	userID, err := strconv.Atoi(request.UserID)
	if err != nil {
		slog.Error("Invalid user ID", "userID", request.UserID, "error", err)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	inverterCreateParams := db.CreateInverterParams{
		UserID:                     int32(userID),
		Vendor:                     request.Vendor,
		Model:                      request.Model,
		SerialNumber:               request.SerialNumber,
		InstallationDate:           request.InstallationDate,
		TotalLifetimeProductionKwh: request.TotalLifetimeProduction,
	}

	inverter, err := uc.inverterQueries.CreateInverter(ctx, inverterCreateParams)
	if err != nil {
		slog.Error("Failed to create inverter in database", "error", err)
		return nil, fmt.Errorf("failed to create inverter: %w", err)
	}

	return &AddInverterResponse{
		ID:                      strconv.FormatInt(int64(inverter.ID), 10),
		UserID:                  strconv.FormatInt(int64(inverter.UserID), 10),
		Vendor:                  inverter.Vendor,
		Model:                   inverter.Model,
		SerialNumber:            inverter.SerialNumber,
		TotalLifetimeProduction: inverter.TotalLifetimeProductionKwh,
		InstallationDate:        inverter.InstallationDate,
	}, nil
}

func (uc *InverterUseCase) LinkInverter(ctx context.Context, userId string, request LinkInverterRequest) (*LinkInverterResponse, error) {
	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	token, err := uc.authClient.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	bearerToken := "Bearer " + token

	resp, err := uc.inverterClient.LinkInverter(timeoutCtx, bearerToken, userId, request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
			slog.Error("link inverter request timed out", "userId", userId, "error", err)
			return nil, fmt.Errorf("link inverter request timed out: %w", err)
		}
		return nil, fmt.Errorf("failed to link inverter: %w", err)
	}

	return resp, nil
}
