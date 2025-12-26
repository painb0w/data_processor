package bootstrap

import (
	"fmt"

	"github.com/painb0w/data_processor/config"
	"github.com/painb0w/data_processor/internal/producer/statsprod"
)

func InitStatsProducer(cfg *config.Config) *statsprod.StatsProducer {

	kafkaBrockers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
	producer := statsprod.NewStatsProducer(kafkaBrockers, cfg.Kafka.StatsUpdateTopic)
	return producer
}
