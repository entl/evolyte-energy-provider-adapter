package enode

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type enodeOAuthResponse struct {
	AccessToken string `json:"access_token" validate:"required"`
	TokenType   string `json:"token_type" validate:"required"`
	ExpiresIn   int    `json:"expires_in" validate:"required"`
	Scope       string `json:"scope" validate:"required"`
}

type enodeAuthClient struct {
	clientID     string
	clientSecret string
	baseURL      string
	oauthBaseURL string
}

func NewEnodeAuthClient(clientID, clientSecret, oauthBaseURL, baseURL string) *enodeAuthClient {
	return &enodeAuthClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		oauthBaseURL: oauthBaseURL,
		baseURL:      baseURL,
	}
}

func (client *enodeAuthClient) Authenticate(c echo.Context) (*enodeOAuthResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	url := fmt.Sprintf("%s/oauth2/token", client.oauthBaseURL)

	slog.Info("Making authentication request", "url", url, "client_id", client.clientID,
		"request_data", form.Encode())

	// Make request
	req, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(client.clientID, client.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{Timeout: 5 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		slog.Error("Failed to make request", "error", err)
		return nil, fmt.Errorf("authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		slog.Error("Authentication failed", "status_code", resp.StatusCode, "response_body", string(bodyBytes))
		return nil, fmt.Errorf("authentication failed: %s", string(bodyBytes))
	}

	// Parse the response body
	var tokenInfo enodeOAuthResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenInfo)
	if err != nil {
		slog.Error("Failed to decode response", "error", err)
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &tokenInfo, nil
}
