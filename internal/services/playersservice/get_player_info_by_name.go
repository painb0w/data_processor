package playersservice

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	pb_models "github.com/painb0w/data_processor/internal/pb/models"
)

func normalizeName(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

func (s *PlayerService) GetPlayerInfoByName(ctx context.Context, name string) ([]*pb_models.PlayerInfo, error) {
	if name == "" {
		return nil, fmt.Errorf("player name cannot be empty")
	}

	normalized := normalizeName(name)

	cached, err := s.cache.GetPlayerInfoByName(ctx, normalized)
	if err == nil && cached != nil {
		return cached, nil
	}

	players, err := s.api.GetPlayerByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch player data from Football API: %w", err)
	}

	if len(players) == 0 {
		return []*pb_models.PlayerInfo{}, nil
	}

	protoPlayers := s.converter.MapPlayerToProto(players)

	cacheTTL := 10 * time.Minute
	if err := s.cache.SetPlayerInfoByName(ctx, normalized, protoPlayers, cacheTTL); err != nil {
		log.Printf("warning: failed to cache player info for name %q: %v", normalized, err)
	}

	return protoPlayers, nil
}
