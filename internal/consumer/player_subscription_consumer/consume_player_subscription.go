package playersubscriptionconsumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type PlayerSubscriptionEvent struct {
	PlayerID int `json:"player_id"`
}

func (c *PlayerSubscriptionConsumer) Consume(ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBroker,
		GroupID:           "data_processor_player_subscriptions_group",
		Topic:             c.topicName,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
	})

	defer reader.Close()

	slog.Info("Starting PlayerSubscriptionConsumer", "topic", c.topicName)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			slog.Error("PlayerSubscriptionConsumer.consume error", "error", err.Error())
		}

		var event PlayerSubscriptionEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			slog.Error("parce", "error", err)
			continue
		}

		err = c.playerSubscriptionsProcessor.Handle(ctx, event.PlayerID)
		if err != nil {
			slog.Error("Handle", "error", err)
		}
	}
}
