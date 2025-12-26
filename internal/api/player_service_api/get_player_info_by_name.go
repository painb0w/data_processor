package player_service_api

import (
	"context"

	pb "github.com/painb0w/data_processor/internal/pb/players_api"
)

func (s *PlayersServiceAPI) GetPlayerInfoByName(ctx context.Context, req *pb.GetPlayerInfoByNameRequest) (*pb.GetPlayerInfoByNameResponse, error) {
	players, err := s.playerService.GetPlayerInfoByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.GetPlayerInfoByNameResponse{PlayerInfos: players}, nil
}
