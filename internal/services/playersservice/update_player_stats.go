package playersservice

import (
	"context"
	"log/slog"
	"reflect"
	"time"

	"github.com/painb0w/data_processor/internal/storage/mongostorage"
)

func (s *PlayerService) snapshotsEqual(a, b *mongostorage.PlayerStats) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Position == b.Position &&
		a.Name == b.Name &&
		reflect.DeepEqual(a.Stats, b.Stats)
}

func (s *PlayerService) UpdateStatsForPlayers(ctx context.Context) error {
	ids, err := s.storage.GetAllTrackedPlayerIDs(ctx)
	if err != nil {
		return err
	}

	for _, id := range ids {
		oldSnap, _ := s.storage.GetPlayerStatsByID(ctx, id)

		playerInfo, err := s.api.GetPlayerByID(ctx, id)
		if err != nil {
			continue
		}

		stats, err := s.api.GetPlayerStatsByPosition(ctx, id, playerInfo.Position)
		if err != nil {
			continue
		}

		newSnap := &mongostorage.PlayerStats{
			ID:        id,
			Name:      playerInfo.Name,
			Position:  playerInfo.Position,
			Stats:     stats,
			UpdatedAt: time.Now(),
		}

		s.storage.UpsertPlayerStats(ctx, newSnap)

		if oldSnap != nil && !s.snapshotsEqual(oldSnap, newSnap) {
			s.producer.ProduceStatsUpdate(ctx, id)
		}
	}
	return nil
}

func (s *PlayerService) StartBackgroundUpdates(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.UpdateStatsForPlayers(ctx); err != nil {
				slog.Error("Failed to update player stats", "error", err)
			}
		case <-ctx.Done():
			slog.Info("Background stats updater stopped")
			return
		}
	}
}
