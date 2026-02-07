package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/s-hammon/p"
)

const MSBaseURL = "https://login.microsoftonline.com/"

type AzureADConfig struct {
	TenantID     string   `json:"tenant_id"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
}

func (cfg AzureADConfig) getScopes() string {
	if len(cfg.Scopes) == 0 {
		return "https://graph.microsoft.com/.default"
	}

	return strings.Join(cfg.Scopes, " ")
}

type TokenResponse struct {
	Token         string `json:"access_token"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	ExtExpiresInt int    `json:"ext_expires_in"`
}

func GetAccessToken(cfg AzureADConfig) (TokenResponse, error) {
	tokenURL := MSBaseURL + p.Format("%s/oauth2/v2.0/token", cfg.TenantID)

	data := url.Values{}
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("scope", cfg.getScopes())
	data.Set("grant_type", "client_credentials")

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode())) // #nosec G107
	if err != nil {
		return TokenResponse{}, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return TokenResponse{}, fmt.Errorf("token request failed with status %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return TokenResponse{}, fmt.Errorf("json.Decode: %v", err)
	}

	return tokenResp, nil
}
