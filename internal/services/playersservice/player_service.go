package playersservice

import (
	"context"
	"time"

	"github.com/painb0w/data_processor/internal/models"
	pb_models "github.com/painb0w/data_processor/internal/pb/models"
	"github.com/painb0w/data_processor/internal/storage/mongostorage"
)

//go:generate mockery --name PlayerCache --output ./mocks --outpkg mocks
type PlayerCache interface {
	GetPlayerInfoByName(ctx context.Context, name string) ([]*pb_models.PlayerInfo, error)
	SetPlayerInfoByName(ctx context.Context, name string, players []*pb_models.PlayerInfo, ttl time.Duration) error
}

//go:generate mockery --name PlayerStorage --output ./mocks --outpkg mocks
type PlayerStorage interface {
	UpsertPlayerStats(ctx context.Context, snapshot *mongostorage.PlayerStats) error
	GetPlayerStatsByID(ctx context.Context, playerID int) (*mongostorage.PlayerStats, error)
	GetAllTrackedPlayerIDs(ctx context.Context) ([]int, error)
}

//go:generate mockery --name StatsProducer --output ./mocks --outpkg mocks
type StatsProducer interface {
	ProduceStatsUpdate(ctx context.Context, playerID int) error
}

type ProtoConverter interface {
	MapPlayerToProto(players []models.FootballAPIPlayerInfo) []*pb_models.PlayerInfo
	MapPlayerStatsToProto(position string, stats any) ([][]byte, error)
}

//go:generate mockery --name FootballAPIClient --output ./mocks --outpkg mocks
type FootballAPIClient interface {
	GetPlayerByName(ctx context.Context, name string) ([]models.FootballAPIPlayerInfo, error)
	GetPlayerByID(ctx context.Context, id int) (*models.FootballAPIPlayerInfo, error)
	GetPlayerStatsByPosition(ctx context.Context, id int, pos string) (any, error)
}

type PlayerService struct {
	api       FootballAPIClient
	cache     PlayerCache
	storage   PlayerStorage
	producer  StatsProducer
	converter ProtoConverter
}

func NewPlayerService(ctx context.Context, api FootballAPIClient, pc PlayerCache, ps PlayerStorage, sp StatsProducer, prc ProtoConverter) *PlayerService {
	return &PlayerService{
		api:       api,
		cache:     pc,
		storage:   ps,
		producer:  sp,
		converter: prc,
	}
}
