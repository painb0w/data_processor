package bootstrap

import (
	"fmt"

	"github.com/painb0w/data_processor/config"

	playersubscriptionsconsumer "github.com/painb0w/data_processor/internal/consumer/player_subscription_consumer"
	playersubscriptionsprocessor "github.com/painb0w/data_processor/internal/services/processors/player_subscriptions_processor"
)

func InitPlayerSubscriptionsConsumer(cfg *config.Config, playerSubscriptionsProcessor *playersubscriptionsprocessor.PlayerSubscriptionsProcessor) *playersubscriptionsconsumer.PlayerSubscriptionConsumer {
	kafkaBrockers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
	return playersubscriptionsconsumer.NewPlayerSubscriptionConsumer(playerSubscriptionsProcessor, kafkaBrockers, cfg.Kafka.SubscriptionTopic)
}
