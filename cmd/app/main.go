package main

import (
	"context"
	"fmt"
	"os"

	"github.com/painb0w/data_processor/config"
	"github.com/painb0w/data_processor/internal/bootstrap"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("error while parcing config, %v", err))
	}

	apiClient := bootstrap.InitFootballApiClient(cfg)
	redisStorage := bootstrap.InitRedisStorage(cfg)
	mongoStorage := bootstrap.InitMongoDb(cfg)
	protoConverter := bootstrap.InitProtoConverter()
	statsProducer := bootstrap.InitStatsProducer(cfg)
	playerService := bootstrap.InitPlayerService(cfg, apiClient, redisStorage, mongoStorage, statsProducer, protoConverter)
	playersApi := bootstrap.InitPlayerServiceAPI(playerService)
	subscriptionsProcessor := bootstrap.InitPlayerSubscriptionsProcessor(playerService)
	subscriptionsConsumer := bootstrap.InitPlayerSubscriptionsConsumer(cfg, subscriptionsProcessor)

	ctx, cancel := context.WithCancel(context.Background())

	bootstrap.AppRun(ctx, cancel, playersApi, subscriptionsConsumer, playerService)
}
