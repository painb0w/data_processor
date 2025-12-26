package player_service_api

import (
	"context"

	pb_models "github.com/painb0w/data_processor/internal/pb/models"
	pb "github.com/painb0w/data_processor/internal/pb/players_api"
)

func (s *PlayersServiceAPI) GetPlayerStatsByID(ctx context.Context, req *pb.GetPlayerStatsByIDRequest) (*pb_models.PlayerStats, error) {
	return s.playerService.GetPlayerStatsByID(ctx, req.Id)
}
