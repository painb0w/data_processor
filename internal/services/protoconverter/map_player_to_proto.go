package protoconverter

import (
	"github.com/painb0w/data_processor/internal/models"
	pb_models "github.com/painb0w/data_processor/internal/pb/models"
)

func (c *ProtoConverter) MapPlayerToProto(players []models.FootballAPIPlayerInfo) []*pb_models.PlayerInfo {
	result := make([]*pb_models.PlayerInfo, len(players))
	for i, p := range players {
		result[i] = &pb_models.PlayerInfo{
			Id:          int32(p.ID),
			Name:        p.Name,
			Age:         int32(p.Age),
			Nationality: p.Nationality,
			Height:      p.Height,
			Weight:      p.Weight,
			Photo:       p.Photo,
			Position:    p.Position,
		}
	}
	return result
}
