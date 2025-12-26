package protoconverter

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/painb0w/data_processor/internal/models"
	pb "github.com/painb0w/data_processor/internal/pb/models"
)

func (c *ProtoConverter) MapPlayerStatsToProto(position string, stats any) ([][]byte, error) {
	var result [][]byte

	switch v := stats.(type) {
	case []models.GoalkeeperStats:
		for _, stat := range v {
			gk := &pb.GoalkeeperStats{
				Gamesplayed:    int32(stat.GamesPlayed),
				Conceded:       int32(stat.Conceded),
				Saves:          int32(stat.Saves),
				Penalties:      int32(stat.Penalties),
				Penaltiessaved: int32(stat.PenaltiesSaved),
			}
			data, err := proto.Marshal(gk)
			if err != nil {
				return nil, err
			}
			result = append(result, data)
		}

	case []models.DefenderStats:
		for _, stat := range v {
			def := &pb.DefenderStats{
				Gamesplayed:   int32(stat.GamesPlayed),
				Tackles:       int32(stat.Tackles),
				Interceptions: int32(stat.Interceptions),
				Blocks:        int32(stat.Blocks),
				Foulscommited: int32(stat.FoulsCommited),
				Duelswon:      int32(stat.DuelsWon),
				Duels:         int32(stat.Duels),
				Yellowcards:   int32(stat.YellowCards),
				Redcards:      int32(stat.RedCards),
			}
			data, err := proto.Marshal(def)
			if err != nil {
				return nil, err
			}
			result = append(result, data)
		}

	case []models.MidfielderStats:
		for _, stat := range v {
			mid := &pb.MidfielderStats{
				Gamesplayed: int32(stat.GamesPlayed),
				Passes:      int32(stat.Passes),
				Keypasses:   int32(stat.KeyPasses),
				Assists:     int32(stat.Assists),
				Tackles:     int32(stat.Tackles),
				Dribbles:    int32(stat.Dribbles),
			}
			data, err := proto.Marshal(mid)
			if err != nil {
				return nil, err
			}
			result = append(result, data)
		}

	case []models.AttackerStats:
		for _, stat := range v {
			att := &pb.AttackerStats{
				Gamesplayed:     int32(stat.GamesPlayed),
				Shots:           int32(stat.Shots),
				Goals:           int32(stat.Goals),
				Assists:         int32(stat.Assists),
				Dribbles:        int32(stat.Dribbles),
				Successdribbles: int32(stat.SuccessDribbles),
				Foulsdrawn:      int32(stat.FoulsDrawn),
			}
			data, err := proto.Marshal(att)
			if err != nil {
				return nil, err
			}
			result = append(result, data)
		}

	default:
		return nil, fmt.Errorf("unsupported stats type for position: %s", position)
	}

	return result, nil
}
