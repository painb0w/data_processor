package playersubscriptionsprocessor

import "context"

func (p *PlayerSubscriptionsProcessor) Handle(ctx context.Context, playerID int) error {
	return p.playerService.SubscribeToPlayerStats(ctx, playerID)
}
