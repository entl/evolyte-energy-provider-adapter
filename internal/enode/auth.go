package enode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type enodeOAuthResponse struct {
	AccessToken string `json:"access_token" validate:"required"`
	TokenType   string `json:"token_type" validate:"required"`
	ExpiresIn   int    `json:"expires_in" validate:"required"`
	Scope       string `json:"scope" validate:"required"`
}

type EnodeAuthClient struct {
	clientID     string
	clientSecret string
	baseURL      string
	oauthBaseURL string
	redisClient  *redis.Client
}

func NewEnodeAuthClient(clientID, clientSecret, oauthBaseURL, baseURL string, redisClient *redis.Client) *EnodeAuthClient {
	return &EnodeAuthClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		oauthBaseURL: oauthBaseURL,
		baseURL:      baseURL,
		redisClient:  redisClient,
	}
}

const (
	enodeAccessTokenKey = "enode_access_token"
)

func (client *EnodeAuthClient) GetAccessToken() (string, error) {
	token, err := client.redisClient.Get(context.Background(), enodeAccessTokenKey).Result()
	if err != redis.Nil {
		slog.Debug("Access token found in Redis")
		return token, nil
	}

	slog.Debug("Access token not found in Redis, authenticating with Enode")
	tokenInfo, err := client.authenticate()
	if err != nil {
		slog.Error("Failed to authenticate with Enode", "error", err)
		return "", err
	}

	if err := client.saveAccessToken(tokenInfo, enodeAccessTokenKey); err != nil {
		slog.Error("Failed to save access token", "error", err)
		return "", err
	}
	return tokenInfo.AccessToken, nil
}

func (client *EnodeAuthClient) authenticate() (*enodeOAuthResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	url := fmt.Sprintf("%s/oauth2/token", client.oauthBaseURL)

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

func (client *EnodeAuthClient) saveAccessToken(tokenInfo *enodeOAuthResponse, key string) error {
	err := client.redisClient.Set(context.Background(), key, tokenInfo.AccessToken, time.Duration(tokenInfo.ExpiresIn)*time.Second-10*time.Second).Err()
	if err != nil {
		slog.Error("Failed to save access token", "error", err)
		return err
	}

	slog.Debug("Saving access token")
	return nil
}
