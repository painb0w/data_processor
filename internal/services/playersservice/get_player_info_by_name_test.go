package playersservice

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/painb0w/data_processor/internal/models"
	pb_models "github.com/painb0w/data_processor/internal/pb/models"
	mocks "github.com/painb0w/data_processor/internal/services/playersservice/mocks"
	"github.com/painb0w/data_processor/internal/services/protoconverter"
	"github.com/stretchr/testify/suite"
)

type GetPlayerInfoByNameSuite struct {
	suite.Suite
	ctx           context.Context
	mockAPI       *mocks.FootballAPIClient
	mockCache     *mocks.PlayerCache
	mockConverter *protoconverter.ProtoConverter
	playerService *PlayerService
}

func (s *GetPlayerInfoByNameSuite) SetupTest() {
	s.ctx = context.Background()
	s.mockAPI = mocks.NewFootballAPIClient(s.T())
	s.mockCache = mocks.NewPlayerCache(s.T())

	mockStorage := mocks.NewPlayerStorage(s.T())
	mockProducer := mocks.NewStatsProducer(s.T())

	s.mockConverter = protoconverter.NewProtoConvert()

	s.playerService = NewPlayerService(
		s.ctx,
		s.mockAPI,
		s.mockCache,
		mockStorage,
		mockProducer,
		s.mockConverter,
	)
}

func (s *GetPlayerInfoByNameSuite) TestGetPlayerInfoByName_ReturnsCachedData() {
	playerName := "Messi"
	normalizedName := "messi"

	cachedPlayers := []*pb_models.PlayerInfo{
		{
			Id:          10,
			Name:        "Messi",
			Position:    "Attacker",
			Nationality: "Argentina",
		},
	}

	s.mockCache.On("GetPlayerInfoByName", s.ctx, normalizedName).Return(cachedPlayers, nil).Once()

	result, err := s.playerService.GetPlayerInfoByName(s.ctx, playerName)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(cachedPlayers, result)
	s.mockCache.AssertExpectations(s.T())
	s.mockAPI.AssertNotCalled(s.T(), "GetPlayerByName")
}

func (s *GetPlayerInfoByNameSuite) TestGetPlayerInfoByName_FallbackToAPIWhenCacheEmpty() {
	playerName := "Ronaldo"
	normalizedName := "ronaldo"

	apiPlayers := []models.FootballAPIPlayerInfo{
		{
			ID:          7,
			Name:        "Ronaldo",
			Position:    "Attacker",
			Age:         38,
			Nationality: "Portugal",
		},
	}

	expectedProtoPlayers := []*pb_models.PlayerInfo{
		{
			Id:          7,
			Name:        "Ronaldo",
			Position:    "Attacker",
			Nationality: "Portugal",
		},
	}

	s.mockCache.On("GetPlayerInfoByName", s.ctx, normalizedName).Return(nil, nil).Once()
	s.mockAPI.On("GetPlayerByName", s.ctx, playerName).Return(apiPlayers, nil).Once()
	s.mockCache.On("SetPlayerInfoByName", s.ctx, normalizedName, expectedProtoPlayers, 10*time.Minute).Return(nil).Once()

	result, err := s.playerService.GetPlayerInfoByName(s.ctx, playerName)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(len(expectedProtoPlayers), len(result))
	s.mockCache.AssertExpectations(s.T())
	s.mockAPI.AssertExpectations(s.T())
}

func (s *GetPlayerInfoByNameSuite) TestGetPlayerInfoByName_NameNormalization() {
	playerName := "  MeSSi  "
	normalizedName := "messi"

	cachedPlayers := []*pb_models.PlayerInfo{
		{Id: 10, Name: "Messi"},
	}

	s.mockCache.On("GetPlayerInfoByName", s.ctx, normalizedName).Return(cachedPlayers, nil).Once()

	result, err := s.playerService.GetPlayerInfoByName(s.ctx, playerName)

	s.NoError(err)
	s.NotNil(result)
	s.mockCache.AssertCalled(s.T(), "GetPlayerInfoByName", s.ctx, normalizedName)
}

func (s *GetPlayerInfoByNameSuite) TestGetPlayerInfoByName_EmptyNameError() {
	playerName := ""

	result, err := s.playerService.GetPlayerInfoByName(s.ctx, playerName)

	s.Error(err)
	s.Nil(result)
	s.Contains(err.Error(), "player name cannot be empty")
	s.mockCache.AssertNotCalled(s.T(), "GetPlayerInfoByName")
	s.mockAPI.AssertNotCalled(s.T(), "GetPlayerByName")
}

func (s *GetPlayerInfoByNameSuite) TestGetPlayerInfoByName_EmptyResponseFromAPI() {
	playerName := "NonExistentPlayer"
	normalizedName := "nonexistentplayer"

	s.mockCache.On("GetPlayerInfoByName", s.ctx, normalizedName).Return(nil, nil).Once()
	s.mockAPI.On("GetPlayerByName", s.ctx, playerName).Return([]models.FootballAPIPlayerInfo{}, nil).Once()

	result, err := s.playerService.GetPlayerInfoByName(s.ctx, playerName)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(0, len(result))
	s.mockCache.AssertNotCalled(s.T(), "SetPlayerInfoByName")
}

func (s *GetPlayerInfoByNameSuite) TestGetPlayerInfoByName_APIError() {
	playerName := "SomePlayer"
	normalizedName := "someplayer"
	expectedErr := errors.New("reached the request limit for the day")

	s.mockCache.On("GetPlayerInfoByName", s.ctx, normalizedName).Return(nil, nil).Once()
	s.mockAPI.On("GetPlayerByName", s.ctx, playerName).Return(nil, expectedErr).Once()

	result, err := s.playerService.GetPlayerInfoByName(s.ctx, playerName)

	s.Error(err)
	s.Nil(result)
	s.Contains(err.Error(), "failed to fetch player data from Football API")
	s.mockAPI.AssertExpectations(s.T())
}

func (s *GetPlayerInfoByNameSuite) TestGetPlayerInfoByName_CacheWriteFailureDoesNotBreak() {
	playerName := "Neymar"
	normalizedName := "neymar"

	apiPlayers := []models.FootballAPIPlayerInfo{
		{
			ID:       11,
			Name:     "Neymar Jr",
			Position: "Attacker",
		},
	}

	expectedProtoPlayers := []*pb_models.PlayerInfo{
		{
			Id:       11,
			Name:     "Neymar Jr",
			Position: "Attacker",
		},
	}

	cacheErr := errors.New("redis connection failed")

	s.mockCache.On("GetPlayerInfoByName", s.ctx, normalizedName).Return(nil, nil).Once()
	s.mockAPI.On("GetPlayerByName", s.ctx, playerName).Return(apiPlayers, nil).Once()
	s.mockCache.On("SetPlayerInfoByName", s.ctx, normalizedName, expectedProtoPlayers, 10*time.Minute).Return(cacheErr).Once()

	result, err := s.playerService.GetPlayerInfoByName(s.ctx, playerName)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(len(expectedProtoPlayers), len(result))
	s.mockCache.AssertExpectations(s.T())
	s.mockAPI.AssertExpectations(s.T())
}

func (s *GetPlayerInfoByNameSuite) TestGetPlayerInfoByName_MultiplePlayersReturned() {
	playerName := "Silva"
	normalizedName := "silva"

	apiPlayers := []models.FootballAPIPlayerInfo{
		{
			ID:       14,
			Name:     "Anderson Silva",
			Position: "Midfielder",
			Age:      33,
		},
		{
			ID:       15,
			Name:     "David Silva",
			Position: "Midfielder",
			Age:      36,
		},
		{
			ID:       16,
			Name:     "Thiago Silva",
			Position: "Defender",
			Age:      38,
		},
	}

	expectedProtoPlayers := []*pb_models.PlayerInfo{
		{Id: 14, Name: "Anderson Silva", Position: "Midfielder"},
		{Id: 15, Name: "David Silva", Position: "Midfielder"},
		{Id: 16, Name: "Thiago Silva", Position: "Defender"},
	}

	s.mockCache.On("GetPlayerInfoByName", s.ctx, normalizedName).Return(nil, nil).Once()
	s.mockAPI.On("GetPlayerByName", s.ctx, playerName).Return(apiPlayers, nil).Once()
	s.mockCache.On("SetPlayerInfoByName", s.ctx, normalizedName, expectedProtoPlayers, 10*time.Minute).Return(nil).Once()

	result, err := s.playerService.GetPlayerInfoByName(s.ctx, playerName)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(3, len(result))
	s.mockCache.AssertExpectations(s.T())
	s.mockAPI.AssertExpectations(s.T())
}

func (s *GetPlayerInfoByNameSuite) TestGetPlayerInfoByName_CacheErrorFallbackToAPI() {
	playerName := "Benzema"
	normalizedName := "benzema"
	cacheErr := errors.New("cache connection timeout")

	apiPlayers := []models.FootballAPIPlayerInfo{
		{
			ID:       9,
			Name:     "Benzema",
			Position: "Attacker",
		},
	}

	expectedProtoPlayers := []*pb_models.PlayerInfo{
		{
			Id:       9,
			Name:     "Benzema",
			Position: "Attacker",
		},
	}

	s.mockCache.On("GetPlayerInfoByName", s.ctx, normalizedName).Return(nil, cacheErr).Once()
	s.mockAPI.On("GetPlayerByName", s.ctx, playerName).Return(apiPlayers, nil).Once()
	s.mockCache.On("SetPlayerInfoByName", s.ctx, normalizedName, expectedProtoPlayers, 10*time.Minute).Return(nil).Once()

	result, err := s.playerService.GetPlayerInfoByName(s.ctx, playerName)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(len(expectedProtoPlayers), len(result))
	s.mockAPI.AssertExpectations(s.T())
}

func (s *GetPlayerInfoByNameSuite) TestGetPlayerInfoByName_WithContextCancellation() {
	playerName := "Haaland"
	cancelCtx, cancel := context.WithCancel(s.ctx)
	cancel()

	result, err := s.playerService.GetPlayerInfoByName(cancelCtx, playerName)

	s.Error(err)
	s.Nil(result)
}

func TestGetPlayerInfoByNameSuite(t *testing.T) {
	suite.Run(t, new(GetPlayerInfoByNameSuite))
}
