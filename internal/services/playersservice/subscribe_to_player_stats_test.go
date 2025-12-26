package playersservice

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/painb0w/data_processor/internal/models"
	mocks "github.com/painb0w/data_processor/internal/services/playersservice/mocks"
	"github.com/painb0w/data_processor/internal/services/protoconverter"
	"github.com/painb0w/data_processor/internal/storage/mongostorage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SubscribeToPlayerStatsSuite struct {
	suite.Suite
	ctx           context.Context
	mockAPI       *mocks.FootballAPIClient
	mockStorage   *mocks.PlayerStorage
	playerService *PlayerService
}

func (s *SubscribeToPlayerStatsSuite) SetupTest() {
	s.ctx = context.Background()
	s.mockAPI = mocks.NewFootballAPIClient(s.T())
	s.mockStorage = mocks.NewPlayerStorage(s.T())

	mockCache := mocks.NewPlayerCache(s.T())
	mockProducer := mocks.NewStatsProducer(s.T())

	s.playerService = NewPlayerService(
		s.ctx,
		s.mockAPI,
		mockCache,
		s.mockStorage,
		mockProducer,
		protoconverter.NewProtoConvert(),
	)
}

func (s *SubscribeToPlayerStatsSuite) TestSubscribeToPlayerStats_SuccessForAttacker() {
	playerID := 100
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:          playerID,
		Name:        "Lewandowski",
		Position:    "Attacker",
		Age:         34,
		Nationality: "Poland",
	}

	stats := models.AttackerStats{
		GamesPlayed:     28,
		Shots:           95,
		Goals:           20,
		Assists:         5,
		Dribbles:        78,
		SuccessDribbles: 52,
		FoulsDrawn:      12,
	}

	s.mockAPI.On("GetPlayerByID", s.ctx, playerID).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, playerID, "Attacker").Return(stats, nil).Once()

	expectedSnapshot := &mongostorage.PlayerStats{
		ID:       playerID,
		Name:     "Lewandowski",
		Position: "Attacker",
		Stats:    stats,
	}

	s.mockStorage.On("UpsertPlayerStats", s.ctx, expectedSnapshot).Return(nil).Once()

	err := s.playerService.SubscribeToPlayerStats(s.ctx, playerID)

	s.NoError(err)
	s.mockAPI.AssertExpectations(s.T())
	s.mockStorage.AssertExpectations(s.T())
}

func (s *SubscribeToPlayerStatsSuite) TestSubscribeToPlayerStats_SuccessForDefender() {
	playerID := 200
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:          playerID,
		Name:        "Gerard Pique",
		Position:    "Defender",
		Age:         35,
		Nationality: "Spain",
	}

	stats := models.DefenderStats{
		GamesPlayed:   22,
		Tackles:       118,
		Interceptions: 42,
		Blocks:        55,
		FoulsCommited: 14,
		DuelsWon:      95,
		Duels:         128,
		YellowCards:   3,
		RedCards:      0,
	}

	s.mockAPI.On("GetPlayerByID", s.ctx, playerID).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, playerID, "Defender").Return(stats, nil).Once()
	s.mockStorage.On("UpsertPlayerStats", s.ctx, expectedSnapshotForStorage(playerID, playerInfo, stats)).Return(nil).Once()

	err := s.playerService.SubscribeToPlayerStats(s.ctx, playerID)

	s.NoError(err)
	s.mockAPI.AssertExpectations(s.T())
}

func (s *SubscribeToPlayerStatsSuite) TestSubscribeToPlayerStats_SuccessForMidfielder() {
	playerID := 300
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:          playerID,
		Name:        "Pedri Gonzalez",
		Position:    "Midfielder",
		Age:         19,
		Nationality: "Spain",
	}

	stats := models.MidfielderStats{
		GamesPlayed: 24,
		Passes:      1650,
		KeyPasses:   32,
		Assists:     4,
		Tackles:     56,
		Dribbles:    98,
	}

	s.mockAPI.On("GetPlayerByID", s.ctx, playerID).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, playerID, "Midfielder").Return(stats, nil).Once()
	s.mockStorage.On("UpsertPlayerStats", s.ctx, expectedSnapshotForStorage(playerID, playerInfo, stats)).Return(nil).Once()

	err := s.playerService.SubscribeToPlayerStats(s.ctx, playerID)

	s.NoError(err)
	s.mockAPI.AssertExpectations(s.T())
}

func (s *SubscribeToPlayerStatsSuite) TestSubscribeToPlayerStats_SuccessForGoalkeeper() {
	playerID := 400
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:          playerID,
		Name:        "Ederson Moraes",
		Position:    "Goalkeeper",
		Age:         30,
		Nationality: "Brazil",
	}

	stats := models.GoalkeeperStats{
		GamesPlayed:    20,
		Conceded:       28,
		Saves:          132,
		Penalties:      3,
		PenaltiesSaved: 1,
	}

	s.mockAPI.On("GetPlayerByID", s.ctx, playerID).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, playerID, "Goalkeeper").Return(stats, nil).Once()
	s.mockStorage.On("UpsertPlayerStats", s.ctx, expectedSnapshotForStorage(playerID, playerInfo, stats)).Return(nil).Once()

	err := s.playerService.SubscribeToPlayerStats(s.ctx, playerID)

	s.NoError(err)
	s.mockAPI.AssertExpectations(s.T())
}

func (s *SubscribeToPlayerStatsSuite) TestSubscribeToPlayerStats_ErrorWhenGetPlayerByIDFails() {
	playerID := 999
	apiErr := errors.New("player not found in external API")

	s.mockAPI.On("GetPlayerByID", s.ctx, playerID).Return(nil, apiErr).Once()

	err := s.playerService.SubscribeToPlayerStats(s.ctx, playerID)

	s.Error(err)
	s.Contains(err.Error(), "failed to get player info")
	s.mockStorage.AssertNotCalled(s.T(), "UpsertPlayerStats")
	s.mockAPI.AssertExpectations(s.T())
}

func (s *SubscribeToPlayerStatsSuite) TestSubscribeToPlayerStats_ErrorWhenGetStatsFails() {
	playerID := 555
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:       playerID,
		Name:     "Test Player",
		Position: "Attacker",
	}
	statsErr := errors.New("stats service temporarily unavailable")

	s.mockAPI.On("GetPlayerByID", s.ctx, playerID).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, playerID, "Attacker").Return(nil, statsErr).Once()

	err := s.playerService.SubscribeToPlayerStats(s.ctx, playerID)

	s.Error(err)
	s.Contains(err.Error(), "failed to get stats")
	s.mockStorage.AssertNotCalled(s.T(), "UpsertPlayerStats")
}

func (s *SubscribeToPlayerStatsSuite) TestSubscribeToPlayerStats_ErrorWhenStorageFails() {
	playerID := 777
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:       playerID,
		Name:     "Storage Test",
		Position: "Defender",
	}
	stats := models.DefenderStats{
		GamesPlayed: 15,
		Tackles:     75,
	}
	storageErr := errors.New("MongoDB connection lost")

	s.mockAPI.On("GetPlayerByID", s.ctx, playerID).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, playerID, "Defender").Return(stats, nil).Once()
	s.mockStorage.On("UpsertPlayerStats", s.ctx, expectedSnapshotForStorage(playerID, playerInfo, stats)).Return(storageErr).Once()

	err := s.playerService.SubscribeToPlayerStats(s.ctx, playerID)

	s.Error(err)
	s.Equal(storageErr, err)
	s.mockStorage.AssertExpectations(s.T())
}

func (s *SubscribeToPlayerStatsSuite) TestSubscribeToPlayerStats_SnapshotContainsCorrectData() {
	playerID := 666
	playerName := "Snapshot Test"
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:       playerID,
		Name:     playerName,
		Position: "Midfielder",
	}
	stats := models.MidfielderStats{
		GamesPlayed: 10,
		Passes:      800,
	}

	s.mockAPI.On("GetPlayerByID", s.ctx, playerID).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, playerID, "Midfielder").Return(stats, nil).Once()

	s.mockStorage.On("UpsertPlayerStats", s.ctx, expectedSnapshotForStorage(playerID, playerInfo, stats)).
		Run(func(args mock.Arguments) {
			snapshot := args.Get(1).(*mongostorage.PlayerStats)
			s.Equal(playerID, snapshot.ID)
			s.Equal(playerName, snapshot.Name)
			s.Equal("Midfielder", snapshot.Position)
			s.Equal(stats, snapshot.Stats)
			s.True(snapshot.UpdatedAt.Before(time.Now().Add(time.Second)))
		}).
		Return(nil).Once()

	err := s.playerService.SubscribeToPlayerStats(s.ctx, playerID)

	s.NoError(err)
}

func (s *SubscribeToPlayerStatsSuite) TestSubscribeToPlayerStats_WithContextCancellation() {
	playerID := 111
	cancelCtx, cancel := context.WithCancel(s.ctx)
	cancel()

	err := s.playerService.SubscribeToPlayerStats(cancelCtx, playerID)

	s.Error(err)
	s.Equal(context.Canceled, err)
}

func (s *SubscribeToPlayerStatsSuite) TestSubscribeToPlayerStats_DifferentPositionsUpsert() {
	for _, pos := range []string{"Attacker", "Defender", "Midfielder", "Goalkeeper"} {
		s.SetupTest()

		playerID := 500
		playerInfo := &models.FootballAPIPlayerInfo{
			ID:       playerID,
			Name:     "Position Test",
			Position: pos,
		}

		stats := getStatsForPosition(pos)

		s.mockAPI.On("GetPlayerByID", s.ctx, playerID).Return(playerInfo, nil).Once()
		s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, playerID, pos).Return(stats, nil).Once()
		s.mockStorage.On("UpsertPlayerStats", s.ctx, expectedSnapshotForStorage(playerID, playerInfo, stats)).Return(nil).Once()

		err := s.playerService.SubscribeToPlayerStats(s.ctx, playerID)

		s.NoError(err)
		s.mockAPI.AssertExpectations(s.T())
		s.mockStorage.AssertExpectations(s.T())
	}
}

func expectedSnapshotForStorage(playerID int, playerInfo *models.FootballAPIPlayerInfo, stats interface{}) *mongostorage.PlayerStats {
	return &mongostorage.PlayerStats{
		ID:       playerID,
		Name:     playerInfo.Name,
		Position: playerInfo.Position,
		Stats:    stats,
	}
}

func getStatsForPosition(position string) interface{} {
	switch position {
	case "Attacker":
		return models.AttackerStats{GamesPlayed: 20, Goals: 10}
	case "Defender":
		return models.DefenderStats{GamesPlayed: 20, Tackles: 100}
	case "Midfielder":
		return models.MidfielderStats{GamesPlayed: 20, Passes: 1500}
	case "Goalkeeper":
		return models.GoalkeeperStats{GamesPlayed: 20, Saves: 120}
	default:
		return nil
	}
}

func TestSubscribeToPlayerStatsSuite(t *testing.T) {
	suite.Run(t, new(SubscribeToPlayerStatsSuite))
}
