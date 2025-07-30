package inverters

import (
	"fmt"

	"github.com/entl/evolyte-energy-provider-adapter/internal/enode"
)

type InverterUseCase struct {
	inverterClient SolarInverterClient
	authClient     *enode.EnodeAuthClient
}

func NewInverterUseCase(inverterClient SolarInverterClient, authClient *enode.EnodeAuthClient) *InverterUseCase {
	return &InverterUseCase{
		inverterClient: inverterClient,
		authClient:     authClient,
	}
}

func (uc *InverterUseCase) ListInverters(after string, before string, pageSize int) (*SolarInverterResponse, error) {
	token, err := uc.authClient.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	bearerToken := "Bearer " + token

	return uc.inverterClient.ListInverters(bearerToken, after, before, pageSize)
}

func (uc *InverterUseCase) ListUserInverters(userID string, after string, before string, pageSize int) (*SolarInverterResponse, error) {
	token, err := uc.authClient.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	bearerToken := "Bearer " + token

	return uc.inverterClient.ListUserInverters(bearerToken, userID, after, before, pageSize)
}

func (uc *InverterUseCase) GetInverter(inverterID string) (*SolarInverter, error) {
	token, err := uc.authClient.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	bearerToken := "Bearer " + token

	inverter, err := uc.inverterClient.GetInverter(bearerToken, inverterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inverter: %w", err)
	}
	return inverter, nil
}

func (uc *InverterUseCase) GetInverterProductionStatistics(inverterID string, year int, month int, day int) (*InverterStatistic, error) {
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

	return uc.inverterClient.GetInverterProductionStatistics(bearerToken, inverterID, params)
}
