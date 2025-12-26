package playersservice

import (
	"context"
	"fmt"
	"time"

	"github.com/painb0w/data_processor/internal/storage/mongostorage"
)

func (s *PlayerService) SubscribeToPlayerStats(ctx context.Context, playerID int) error {
	playerInfo, err := s.api.GetPlayerByID(ctx, playerID)
	if err != nil {
		return fmt.Errorf("failed to get player info: %w", err)
	}

	stats, err := s.api.GetPlayerStatsByPosition(ctx, playerID, playerInfo.Position)
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	snapshot := &mongostorage.PlayerStats{
		ID:        playerID,
		Name:      playerInfo.Name,
		Position:  playerInfo.Position,
		Stats:     stats,
		UpdatedAt: time.Now(),
	}

	return s.storage.UpsertPlayerStats(ctx, snapshot)
}
