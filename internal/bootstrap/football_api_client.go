package bootstrap

import (
	"github.com/painb0w/data_processor/config"
	"github.com/painb0w/data_processor/internal/services/footballapi"
)

func InitFootballApiClient(cfg *config.Config) *footballapi.FootballApiClient {
	client := footballapi.NewFootballApiClient(cfg.FootballApi.BaseUrl, cfg.FootballApi.ApiKey, cfg.FootballApi.Timeout)
	return client
}
