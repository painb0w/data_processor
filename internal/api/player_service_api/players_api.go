package player_service_api

import (
	"context"

	pb_models "github.com/painb0w/data_processor/internal/pb/models"
	pb "github.com/painb0w/data_processor/internal/pb/players_api"
)

type playerService interface {
	GetPlayerInfoByName(ctx context.Context, name string) ([]*pb_models.PlayerInfo, error)
	GetPlayerStatsByID(ctx context.Context, id int32) (*pb_models.PlayerStats, error)
}

type PlayersServiceAPI struct {
	pb.UnimplementedPlayersServiceServer
	playerService playerService
}

func NewPlayersServiceAPI(s playerService) *PlayersServiceAPI {
	return &PlayersServiceAPI{
		playerService: s,
	}
}
