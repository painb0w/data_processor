package redisstorage

import (
	"context"

	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"

	pb_models "github.com/painb0w/data_processor/internal/pb/models"
	players_api "github.com/painb0w/data_processor/internal/pb/players_api"
)

func (r *RedisStorage) GetPlayerInfoByName(ctx context.Context, name string) ([]*pb_models.PlayerInfo, error) {
	key := "player:name:" + name
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var resp players_api.GetPlayerInfoByNameResponse
	if err := proto.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.PlayerInfos, nil
}
