package playersservice

import (
	"context"
	"fmt"
	"time"

	pb_models "github.com/painb0w/data_processor/internal/pb/models"
)

func (s *PlayerService) GetPlayerStatsByID(ctx context.Context, id int32) (*pb_models.PlayerStats, error) {
	playerInfo, err := s.api.GetPlayerByID(ctx, int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get player info for ID %d: %w", id, err)
	}

	stats, err := s.api.GetPlayerStatsByPosition(ctx, int(id), playerInfo.Position)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats for player ID %d: %w", id, err)
	}

	protoStatsList, err := s.converter.MapPlayerStatsToProto(playerInfo.Position, stats)
	if err != nil {
		return nil, fmt.Errorf("failed to convert stats to protobuf for player ID %d: %w", id, err)
	}

	return &pb_models.PlayerStats{
		Id:        id,
		Name:      playerInfo.Name,
		Position:  playerInfo.Position,
		Stats:     protoStatsList,
		UpdatedAt: time.Now().Unix(),
	}, nil
}
