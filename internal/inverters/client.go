package inverters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/entl/evolyte-energy-provider-adapter/internal/db"
	"github.com/entl/evolyte-energy-provider-adapter/internal/enode"
)

type SolarInverterClient interface {
	ListInverters(ctx context.Context, bearbearerToken string, after string, before string, pageSize int) (*SolarInverterResponse, error)
	ListUserInverters(ctx context.Context, bearerToken string, userID string, after string, before string, pageSize int) (*SolarInverterResponse, error)
	GetInverter(ctx context.Context, bearerToken string, inverterID string) (*SolarInverter, error)
	GetInverterProductionStatistics(ctx context.Context, bearerToken string, inverterID string, params InverterStatisticParams) (*InverterStatistic, error)
	LinkInverter(ctx context.Context, bearerToken string, userId string, linkBody LinkInverterRequest) (*LinkInverterResponse, error)
}

type EnodeSolarInverterClient struct {
	enodeAuthClient *enode.EnodeAuthClient
	enodeBaseURL    string
	inverterQueries *db.Queries
	httpClient      *http.Client
}

func NewEnodeSolarInverterClient(authClient *enode.EnodeAuthClient, baseURL string, httpClient *http.Client, inverterQueries *db.Queries) *EnodeSolarInverterClient {
	return &EnodeSolarInverterClient{
		enodeAuthClient: authClient,
		enodeBaseURL:    baseURL,
		inverterQueries: inverterQueries,
		httpClient:      httpClient,
	}
}

func (client *EnodeSolarInverterClient) ListInverters(ctx context.Context, bearerToken string, after string, before string, pageSize int) (*SolarInverterResponse, error) {
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

	if response.StatusCode != http.StatusOK {
		return nil, parseEnodeError(response)
	}

	var solarInverterResponse SolarInverterResponse
	if err := json.NewDecoder(response.Body).Decode(&solarInverterResponse); err != nil {
		return nil, err
	}

	return &solarInverterResponse, nil
}

func (client *EnodeSolarInverterClient) ListUserInverters(ctx context.Context, bearerToken string, userID string, after string, before string, pageSize int) (*SolarInverterResponse, error) {
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

func (client *EnodeSolarInverterClient) GetInverter(ctx context.Context, bearerToken string, inverterID string) (*SolarInverter, error) {
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
		return nil, parseEnodeError(response)
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

func (client *EnodeSolarInverterClient) GetInverterProductionStatistics(ctx context.Context, bearerToken string, inverterID string, inverterStatisticParams InverterStatisticParams) (*InverterStatistic, error) {
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

	if response.StatusCode != http.StatusOK {
		return nil, parseEnodeError(response)
	}

	var inverterStatistic InverterStatistic
	if err := json.NewDecoder(response.Body).Decode(&inverterStatistic); err != nil {
		return nil, err
	}

	return &inverterStatistic, nil
}

func (client *EnodeSolarInverterClient) LinkInverter(ctx context.Context, bearerToken string, userId string, linkBody LinkInverterRequest) (*LinkInverterResponse, error) {
	linkURL := fmt.Sprintf("%s/users/%s/link", client.enodeBaseURL, userId)
	reqBody, err := json.Marshal(linkBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", linkURL, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("Content-Type", "application/json")

	response, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, parseEnodeError(response)
	}

	var linkResponse LinkInverterResponse
	if err := json.NewDecoder(response.Body).Decode(&linkResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &linkResponse, nil
}

func parseEnodeError(resp *http.Response) error {
	defer resp.Body.Close()

	var apiErr struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		return fmt.Errorf("request failed with status %s and unreadable error body: %w", resp.Status, err)
	}

	return fmt.Errorf("%s - %s", apiErr.Title, apiErr.Detail)
}
