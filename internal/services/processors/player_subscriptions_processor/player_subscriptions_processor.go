package playersubscriptionsprocessor

import (
	"context"
)

type playerService interface {
	SubscribeToPlayerStats(ctx context.Context, playerID int) error
}

type PlayerSubscriptionsProcessor struct {
	playerService playerService
}

func NewPlayerSubscriptionsProcessor(s playerService) *PlayerSubscriptionsProcessor {
	return &PlayerSubscriptionsProcessor{
		playerService: s,
	}
}
