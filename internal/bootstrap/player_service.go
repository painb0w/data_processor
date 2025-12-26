package bootstrap

import (
	"context"

	"github.com/painb0w/data_processor/config"
	"github.com/painb0w/data_processor/internal/producer/statsprod"
	footballapi "github.com/painb0w/data_processor/internal/services/footballapi"
	"github.com/painb0w/data_processor/internal/services/playersservice"
	"github.com/painb0w/data_processor/internal/services/protoconverter"
	"github.com/painb0w/data_processor/internal/storage/mongostorage"
	"github.com/painb0w/data_processor/internal/storage/redisstorage"
)

func InitPlayerService(cfg *config.Config, api *footballapi.FootballApiClient, cache *redisstorage.RedisStorage, ms *mongostorage.MongoStorage, sp *statsprod.StatsProducer, pc *protoconverter.ProtoConverter) *playersservice.PlayerService {
	return playersservice.NewPlayerService(context.Background(), api, cache, ms, sp, pc)
}
