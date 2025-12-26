package playersservice

import (
	"context"
	"testing"

	mocks "github.com/painb0w/data_processor/internal/services/playersservice/mocks"
	"github.com/painb0w/data_processor/internal/services/protoconverter"
	"github.com/stretchr/testify/suite"
)

type PlayerServiceSuite struct {
	suite.Suite
	ctx            context.Context
	mockAPI        *mocks.FootballAPIClient
	mockMongo      *mocks.PlayerStorage
	mockProducer   *mocks.StatsProducer
	mockCache      *mocks.PlayerCache
	playerService  *PlayerService
	protoConverter *ProtoConverter
}

func (s *PlayerServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.mockAPI = mocks.NewFootballAPIClient(s.T())
	s.mockMongo = mocks.NewPlayerStorage(s.T())
	s.mockCache = mocks.NewPlayerCache(s.T())
	s.mockProducer = mocks.NewStatsProducer(s.T())

	protoConverter := protoconverter.NewProtoConvert()

	s.playerService = NewPlayerService(
		s.ctx,
		s.mockAPI,
		s.mockCache,
		s.mockMongo,
		s.mockProducer,
		protoConverter,
	)
}

func TestPlayerServiceSuite(t *testing.T) {
	suite.Run(t, new(PlayerServiceSuite))
}
