package playersservice

import (
	"context"
	"errors"
	"testing"

	"github.com/painb0w/data_processor/internal/models"
	mocks "github.com/painb0w/data_processor/internal/services/playersservice/mocks"
	"github.com/painb0w/data_processor/internal/services/protoconverter"
	"github.com/stretchr/testify/suite"
)

type GetPlayerStatsByIDSuite struct {
	suite.Suite
	ctx           context.Context
	mockAPI       *mocks.FootballAPIClient
	mockConverter *protoconverter.ProtoConverter
	playerService *PlayerService
}

func (s *GetPlayerStatsByIDSuite) SetupTest() {
	s.ctx = context.Background()
	s.mockAPI = mocks.NewFootballAPIClient(s.T())

	mockCache := mocks.NewPlayerCache(s.T())
	mockStorage := mocks.NewPlayerStorage(s.T())
	mockProducer := mocks.NewStatsProducer(s.T())

	s.mockConverter = protoconverter.NewProtoConvert()

	s.playerService = NewPlayerService(
		s.ctx,
		s.mockAPI,
		mockCache,
		mockStorage,
		mockProducer,
		s.mockConverter,
	)
}

func (s *GetPlayerStatsByIDSuite) TestGetPlayerStatsByID_SuccessForAttacker() {
	playerID := int32(123)
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:          int(playerID),
		Name:        "Ronaldo",
		Position:    "Attacker",
		Age:         38,
		Nationality: "Portugal",
		Height:      "6.18",
		Weight:      "84",
	}

	stats := models.AttackerStats{
		GamesPlayed:     30,
		Shots:           120,
		Goals:           25,
		Assists:         8,
		Dribbles:        145,
		SuccessDribbles: 95,
		FoulsDrawn:      22,
	}

	s.mockAPI.On("GetPlayerByID", s.ctx, int(playerID)).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, int(playerID), "Attacker").Return(stats, nil).Once()

	result, err := s.playerService.GetPlayerStatsByID(s.ctx, playerID)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(playerID, result.Id)
	s.Equal("Ronaldo", result.Name)
	s.Equal("Attacker", result.Position)
	s.Equal(1, len(result.Stats))
	s.True(result.UpdatedAt > 0)
	s.mockAPI.AssertExpectations(s.T())
}

func (s *GetPlayerStatsByIDSuite) TestGetPlayerStatsByID_SuccessForDefender() {
	playerID := int32(456)
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:          int(playerID),
		Name:        "Ramos",
		Position:    "Defender",
		Age:         36,
		Nationality: "Spain",
		Height:      "6.0",
		Weight:      "82",
	}

	stats := models.DefenderStats{
		GamesPlayed:   28,
		Tackles:       145,
		Interceptions: 52,
		Blocks:        68,
		FoulsCommited: 18,
		DuelsWon:      112,
		Duels:         156,
		YellowCards:   5,
		RedCards:      0,
	}

	s.mockAPI.On("GetPlayerByID", s.ctx, int(playerID)).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, int(playerID), "Defender").Return(stats, nil).Once()

	result, err := s.playerService.GetPlayerStatsByID(s.ctx, playerID)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(playerID, result.Id)
	s.Equal("Ramos", result.Name)
	s.Equal("Defender", result.Position)
	s.mockAPI.AssertExpectations(s.T())
}

func (s *GetPlayerStatsByIDSuite) TestGetPlayerStatsByID_SuccessForMidfielder() {
	playerID := int32(789)
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:          int(playerID),
		Name:        "Modric",
		Position:    "Midfielder",
		Age:         37,
		Nationality: "Croatia",
		Height:      "5.82",
		Weight:      "75",
	}

	stats := models.MidfielderStats{
		GamesPlayed: 26,
		Passes:      1850,
		KeyPasses:   45,
		Assists:     6,
		Tackles:     78,
		Dribbles:    124,
	}

	s.mockAPI.On("GetPlayerByID", s.ctx, int(playerID)).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, int(playerID), "Midfielder").Return(stats, nil).Once()

	result, err := s.playerService.GetPlayerStatsByID(s.ctx, playerID)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(playerID, result.Id)
	s.Equal("Modric", result.Name)
	s.Equal("Midfielder", result.Position)
	s.mockAPI.AssertExpectations(s.T())
}

func (s *GetPlayerStatsByIDSuite) TestGetPlayerStatsByID_SuccessForGoalkeeper() {
	playerID := int32(101112)
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:          int(playerID),
		Name:        "Neuer",
		Position:    "Goalkeeper",
		Age:         37,
		Nationality: "Germany",
		Height:      "6.3",
		Weight:      "92",
	}

	stats := models.GoalkeeperStats{
		GamesPlayed:    25,
		Conceded:       32,
		Saves:          145,
		Penalties:      4,
		PenaltiesSaved: 1,
	}

	s.mockAPI.On("GetPlayerByID", s.ctx, int(playerID)).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, int(playerID), "Goalkeeper").Return(stats, nil).Once()

	result, err := s.playerService.GetPlayerStatsByID(s.ctx, playerID)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(playerID, result.Id)
	s.Equal("Neuer", result.Name)
	s.Equal("Goalkeeper", result.Position)
	s.mockAPI.AssertExpectations(s.T())
}

func (s *GetPlayerStatsByIDSuite) TestGetPlayerStatsByID_ErrorWhenGetPlayerByIDFails() {
	playerID := int32(999)
	expectedErr := errors.New("API error: player not found")

	s.mockAPI.On("GetPlayerByID", s.ctx, int(playerID)).Return(nil, expectedErr).Once()

	result, err := s.playerService.GetPlayerStatsByID(s.ctx, playerID)

	s.Error(err)
	s.Nil(result)
	s.Contains(err.Error(), "failed to get player info")
	s.mockAPI.AssertExpectations(s.T())
}

func (s *GetPlayerStatsByIDSuite) TestGetPlayerStatsByID_ErrorWhenGetPlayerStatsFails() {
	playerID := int32(555)
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:       int(playerID),
		Name:     "Test Player",
		Position: "Attacker",
	}
	expectedErr := errors.New("API error: stats service unavailable")

	s.mockAPI.On("GetPlayerByID", s.ctx, int(playerID)).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, int(playerID), "Attacker").Return(nil, expectedErr).Once()

	result, err := s.playerService.GetPlayerStatsByID(s.ctx, playerID)

	s.Error(err)
	s.Nil(result)
	s.Contains(err.Error(), "failed to get stats for player")
	s.mockAPI.AssertExpectations(s.T())
}

func (s *GetPlayerStatsByIDSuite) TestGetPlayerStatsByID_ErrorWhenConversionFails() {
	playerID := int32(777)
	playerInfo := &models.FootballAPIPlayerInfo{
		ID:       int(playerID),
		Name:     "Test Player",
		Position: "UnknownPosition",
	}

	stats := models.AttackerStats{
		GamesPlayed: 20,
	}

	s.mockAPI.On("GetPlayerByID", s.ctx, int(playerID)).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, int(playerID), "UnknownPosition").Return(stats, nil).Once()

	result, err := s.playerService.GetPlayerStatsByID(s.ctx, playerID)

	s.Error(err)
	s.Nil(result)
	s.Contains(err.Error(), "failed to convert stats to protobuf")
	s.mockAPI.AssertExpectations(s.T())
}

func TestGetPlayerStatsByIDSuite(t *testing.T) {
	suite.Run(t, new(GetPlayerStatsByIDSuite))
}
