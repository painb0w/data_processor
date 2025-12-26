package footballapi

import (
	"net/http"
	"time"
)

type FootballApiClient struct {
	BaseURL string
	ApiKey  string
	Client  *http.Client
}

func NewFootballApiClient(baseURL, apiKey string, timeoutSeconds int) *FootballApiClient {
	return &FootballApiClient{
		BaseURL: baseURL,
		ApiKey:  apiKey,
		Client:  &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second},
	}
}
