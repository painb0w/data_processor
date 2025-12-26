package redisstorage

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	pb_models "github.com/painb0w/data_processor/internal/pb/models"
	players_api "github.com/painb0w/data_processor/internal/pb/players_api"
)

func (r *RedisStorage) SetPlayerInfoByName(ctx context.Context, name string, players []*pb_models.PlayerInfo, ttl time.Duration) error {
	key := " player:name:" + name

	resp := &players_api.GetPlayerInfoByNameResponse{
		PlayerInfos: players,
	}
	data, err := proto.Marshal(resp)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}
