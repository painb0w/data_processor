package statsprod

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type StatsProducer struct {
	writer      *kafka.Writer
	kafkaBroker []string
	topicName   string
}

func NewStatsProducer(kafkaBroker []string, topicName string) *StatsProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaBroker...),
		Topic:        topicName,
		Balancer:     &kafka.Hash{},
		BatchTimeout: 10 * time.Millisecond,
		BatchSize:    100,
		RequiredAcks: kafka.RequireAll,
		Async:        false,
		Compression:  kafka.Snappy,
	}

	return &StatsProducer{
		writer:      writer,
		kafkaBroker: kafkaBroker,
		topicName:   topicName,
	}
}
