package playersubscriptionconsumer

import (
	"context"
)

type playerSubscriptionsProcessor interface {
	Handle(ctx context.Context, playerID int) error
}

type PlayerSubscriptionConsumer struct {
	playerSubscriptionsProcessor playerSubscriptionsProcessor
	kafkaBroker                  []string
	topicName                    string
}

func NewPlayerSubscriptionConsumer(psp playerSubscriptionsProcessor, kafkaBroker []string, topicName string) *PlayerSubscriptionConsumer {
	return &PlayerSubscriptionConsumer{
		playerSubscriptionsProcessor: psp,
		kafkaBroker:                  kafkaBroker,
		topicName:                    topicName,
	}
}
