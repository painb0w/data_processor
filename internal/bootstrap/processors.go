package bootstrap

import (
	"github.com/painb0w/data_processor/internal/services/playersservice"
	playersubscriptionsprocessor "github.com/painb0w/data_processor/internal/services/processors/player_subscriptions_processor"
)

func InitPlayerSubscriptionsProcessor(playerService *playersservice.PlayerService) *playersubscriptionsprocessor.PlayerSubscriptionsProcessor {
	return playersubscriptionsprocessor.NewPlayerSubscriptionsProcessor(playerService)
}
