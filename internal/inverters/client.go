package inverters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/entl/evolyte-energy-provider-adapter/internal/enode"
)

type SolarInverterClient interface {
	ListInverters(bearbearerToken string, after string, before string, pageSize int) (*SolarInverterResponse, error)
	ListUserInverters(bearerToken string, userID string, after string, before string, pageSize int) (*SolarInverterResponse, error)
	GetInverter(bearerToken string, inverterID string) (*SolarInverter, error)
	GetInverterProductionStatistics(bearerToken string, inverterID string, params InverterStatisticParams) (*InverterStatistic, error)
}

type EnodeSolarInverterClient struct {
	enodeAuthClient *enode.EnodeAuthClient
	enodeBaseURL    string
	httpClient      *http.Client
}

func NewEnodeSolarInverterClient(authClient *enode.EnodeAuthClient, baseURL string, httpClient *http.Client) *EnodeSolarInverterClient {
	return &EnodeSolarInverterClient{
		enodeAuthClient: authClient,
		enodeBaseURL:    baseURL,
		httpClient:      httpClient,
	}
}

func (client *EnodeSolarInverterClient) ListInverters(bearerToken string, after string, before string, pageSize int) (*SolarInverterResponse, error) {
	invertersBaseURL := client.enodeBaseURL + "/inverters"
	token, err := client.enodeAuthClient.GetAccessToken()
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	// Add pagination parameters
	client.addPaginationParams(&params, after, before, pageSize)

	fullURL := invertersBaseURL
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		return nil, err
	}

	response, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var solarInverterResponse SolarInverterResponse
	if err := json.NewDecoder(response.Body).Decode(&solarInverterResponse); err != nil {
		return nil, err
	}

	return &solarInverterResponse, nil
}

func (client *EnodeSolarInverterClient) ListUserInverters(bearerToken string, userID string, after string, before string, pageSize int) (*SolarInverterResponse, error) {
	invertersBaseURL := fmt.Sprintf("%s/users/%s/inverters", client.enodeBaseURL, userID)
	token, err := client.enodeAuthClient.GetAccessToken()
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	// Add pagination parameters
	client.addPaginationParams(&params, after, before, pageSize)

	fullURL := invertersBaseURL
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}
	req, err := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		return nil, err
	}
	response, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var solarInverterResponse SolarInverterResponse
	if err := json.NewDecoder(response.Body).Decode(&solarInverterResponse); err != nil {
		return nil, err
	}

	return &solarInverterResponse, nil
}

func (client *EnodeSolarInverterClient) addPaginationParams(params *url.Values, after string, before string, pageSize int) *url.Values {
	if after != "" {
		params.Add("after", after)
	}
	if before != "" {
		params.Add("before", before)
	}
	if pageSize > 0 {
		params.Add("pageSize", strconv.Itoa(pageSize))
	}
	return params
}

func (client *EnodeSolarInverterClient) GetInverter(bearerToken string, inverterID string) (*SolarInverter, error) {
	inverterURL := fmt.Sprintf("%s/inverters/%s", client.enodeBaseURL, inverterID)
	token, err := client.enodeAuthClient.GetAccessToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", inverterURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	response, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get inverter: %s", response.Status)
	}

	var inverter SolarInverter
	if err := json.NewDecoder(response.Body).Decode(&inverter); err != nil {
		return nil, err
	}

	return &inverter, nil
}

type InverterStatisticParams struct {
	Year  int `validate:"required"`
	Month int `validate:"required"`
	Day   int
}

func (p InverterStatisticParams) Validate() error {
	if p.Year <= 0 {
		return fmt.Errorf("invalid year: %d", p.Year)
	}
	if p.Month < 1 || p.Month > 12 {
		return fmt.Errorf("invalid month: %d", p.Month)
	}
	if p.Day < 0 || p.Day > 31 {
		return fmt.Errorf("invalid day: %d", p.Day)
	}
	return nil
}

func (client *EnodeSolarInverterClient) GetInverterProductionStatistics(bearerToken string, inverterID string, inverterStatisticParams InverterStatisticParams) (*InverterStatistic, error) {
	inverterURL := fmt.Sprintf("%s/inverters/%s/statistics", client.enodeBaseURL, inverterID)
	if err := inverterStatisticParams.Validate(); err != nil {
		return nil, fmt.Errorf("invalid inverter statistic parameters: %w", err)
	}

	params := url.Values{}
	params.Add("year", strconv.Itoa(inverterStatisticParams.Year))
	params.Add("month", strconv.Itoa(inverterStatisticParams.Month))
	if inverterStatisticParams.Day > 0 {
		params.Add("day", strconv.Itoa(inverterStatisticParams.Day))
	}

	fullURL := inverterURL + "?" + params.Encode()
	req, err := http.NewRequest("GET", fullURL, nil)
	req.Header.Set("Authorization", bearerToken)
	if err != nil {
		return nil, err
	}
	response, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var inverterStatistic InverterStatistic
	if err := json.NewDecoder(response.Body).Decode(&inverterStatistic); err != nil {
		return nil, err
	}

	return &inverterStatistic, nil
}
