package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

func GetUrlGoogleString() (string, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		return "", fmt.Errorf("GOOGLE_CLIENT_ID environment variable not found")
	}

	redirectURI := os.Getenv("GOOGLE_URI_REDIRECT")
	if redirectURI == "" {
		return "", fmt.Errorf("GOOGLE_URI_REDIRECT environment variable not found")
	}

	baseURL := "https://accounts.google.com/o/oauth2/v2/auth"
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %v", err)
	}

	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", redirectURI)
	params.Add("scope", "openid email profile")
	params.Add("access_type", "offline")

	u.RawQuery = params.Encode()

	return u.String(), nil
}

func ExchangeGoogleCodeForToken(ctx context.Context, code string) (*GoogleTokenResponse, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID environment variable not found")
	}

	clientSecret := os.Getenv("GOOGLE_SECRET")
	if clientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_SECRET environment variable not found")
	}

	redirectURI := os.Getenv("GOOGLE_URI_REDIRECT")
	if redirectURI == "" {
		return nil, fmt.Errorf("GOOGLE_URI_REDIRECT environment variable not found")
	}

	// Prepare form data
	params := url.Values{}
	params.Set("code", code)
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)
	params.Set("redirect_uri", redirectURI)
	params.Set("grant_type", "authorization_code")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tokenResponse GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	return &tokenResponse, nil
}
