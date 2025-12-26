package playersservice

import (
	"context"
	"errors"
	"testing"

	"github.com/painb0w/data_processor/internal/models"
	mocks "github.com/painb0w/data_processor/internal/services/playersservice/mocks"
	"github.com/painb0w/data_processor/internal/services/protoconverter"
	"github.com/painb0w/data_processor/internal/storage/mongostorage"
	"github.com/stretchr/testify/suite"
)

type UpdatePlayerStatsSuite struct {
	suite.Suite
	ctx           context.Context
	mockAPI       *mocks.FootballAPIClient
	mockStorage   *mocks.PlayerStorage
	mockProducer  *mocks.StatsProducer
	playerService *PlayerService
}

func (s *UpdatePlayerStatsSuite) SetupTest() {
	s.ctx = context.Background()
	s.mockAPI = mocks.NewFootballAPIClient(s.T())
	s.mockStorage = mocks.NewPlayerStorage(s.T())
	s.mockProducer = mocks.NewStatsProducer(s.T())

	mockCache := mocks.NewPlayerCache(s.T())

	s.playerService = NewPlayerService(
		s.ctx,
		s.mockAPI,
		mockCache,
		s.mockStorage,
		s.mockProducer,
		protoconverter.NewProtoConvert(),
	)
}

func (s *UpdatePlayerStatsSuite) TestUpdateStatsForPlayers_SuccessfullyUpdatesAllPlayers() {
	playerIDs := []int{1, 2, 3}

	s.mockStorage.On("GetAllTrackedPlayerIDs", s.ctx).Return(playerIDs, nil).Once()

	for _, id := range playerIDs {
		playerInfo := &models.FootballAPIPlayerInfo{
			ID:       id,
			Name:     "Player " + string(rune(id+'0')),
			Position: "Attacker",
		}
		stats := models.AttackerStats{GamesPlayed: 10, Goals: 5}

		s.mockStorage.On("GetPlayerStatsByID", s.ctx, id).Return(nil, nil).Once()
		s.mockAPI.On("GetPlayerByID", s.ctx, id).Return(playerInfo, nil).Once()
		s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, id, "Attacker").Return(stats, nil).Once()
		s.mockStorage.On("UpsertPlayerStats", s.ctx, matchPlayerStats(id, "Player "+string(rune(id+'0')))).Return(nil).Once()
	}

	err := s.playerService.UpdateStatsForPlayers(s.ctx)

	s.NoError(err)
	s.mockStorage.AssertExpectations(s.T())
	s.mockAPI.AssertExpectations(s.T())
}

func (s *UpdatePlayerStatsSuite) TestUpdateStatsForPlayers_SkipsPlayersWithAPIErrors() {
	playerIDs := []int{10, 20, 30}

	s.mockStorage.On("GetAllTrackedPlayerIDs", s.ctx).Return(playerIDs, nil).Once()

	playerInfo1 := &models.FootballAPIPlayerInfo{ID: 10, Name: "Player1", Position: "Attacker"}
	stats1 := models.AttackerStats{GamesPlayed: 10}
	s.mockStorage.On("GetPlayerStatsByID", s.ctx, 10).Return(nil, nil).Once()
	s.mockAPI.On("GetPlayerByID", s.ctx, 10).Return(playerInfo1, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, 10, "Attacker").Return(stats1, nil).Once()
	s.mockStorage.On("UpsertPlayerStats", s.ctx, matchPlayerStats(10, "Player1")).Return(nil).Once()

	s.mockStorage.On("GetPlayerStatsByID", s.ctx, 20).Return(nil, nil).Once()
	s.mockAPI.On("GetPlayerByID", s.ctx, 20).Return(nil, errors.New("API error")).Once()

	playerInfo3 := &models.FootballAPIPlayerInfo{ID: 30, Name: "Player3", Position: "Defender"}
	s.mockStorage.On("GetPlayerStatsByID", s.ctx, 30).Return(nil, nil).Once()
	s.mockAPI.On("GetPlayerByID", s.ctx, 30).Return(playerInfo3, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, 30, "Defender").Return(nil, errors.New("stats error")).Once()

	err := s.playerService.UpdateStatsForPlayers(s.ctx)

	s.NoError(err)
	s.mockAPI.AssertExpectations(s.T())
}

func (s *UpdatePlayerStatsSuite) TestUpdateStatsForPlayers_ProducesUpdateWhenStatsChanged() {
	playerID := 100

	oldStats := models.AttackerStats{GamesPlayed: 10, Goals: 5, Assists: 2}
	oldSnap := &mongostorage.PlayerStats{
		ID:       playerID,
		Name:     "Changed Player",
		Position: "Attacker",
		Stats:    oldStats,
	}

	s.mockStorage.On("GetAllTrackedPlayerIDs", s.ctx).Return([]int{playerID}, nil).Once()
	s.mockStorage.On("GetPlayerStatsByID", s.ctx, playerID).Return(oldSnap, nil).Once()

	playerInfo := &models.FootballAPIPlayerInfo{ID: playerID, Name: "Changed Player", Position: "Attacker"}
	newStats := models.AttackerStats{GamesPlayed: 11, Goals: 6, Assists: 2} // Изменилась статистика

	s.mockAPI.On("GetPlayerByID", s.ctx, playerID).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, playerID, "Attacker").Return(newStats, nil).Once()
	s.mockStorage.On("UpsertPlayerStats", s.ctx, matchPlayerStats(playerID, "Changed Player")).Return(nil).Once()

	s.mockProducer.On("ProduceStatsUpdate", s.ctx, playerID).Return(nil).Once()

	err := s.playerService.UpdateStatsForPlayers(s.ctx)

	s.NoError(err)
	s.mockProducer.AssertExpectations(s.T())
}

func (s *UpdatePlayerStatsSuite) TestUpdateStatsForPlayers_DoesNotProduceUpdateWhenStatsUnchanged() {
	playerID := 200

	stats := models.AttackerStats{GamesPlayed: 10, Goals: 5, Assists: 2}
	oldSnap := &mongostorage.PlayerStats{
		ID:       playerID,
		Name:     "Unchanged Player",
		Position: "Attacker",
		Stats:    stats,
	}

	s.mockStorage.On("GetAllTrackedPlayerIDs", s.ctx).Return([]int{playerID}, nil).Once()
	s.mockStorage.On("GetPlayerStatsByID", s.ctx, playerID).Return(oldSnap, nil).Once()

	playerInfo := &models.FootballAPIPlayerInfo{ID: playerID, Name: "Unchanged Player", Position: "Attacker"}

	s.mockAPI.On("GetPlayerByID", s.ctx, playerID).Return(playerInfo, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, playerID, "Attacker").Return(stats, nil).Once()
	s.mockStorage.On("UpsertPlayerStats", s.ctx, matchPlayerStats(playerID, "Unchanged Player")).Return(nil).Once()

	err := s.playerService.UpdateStatsForPlayers(s.ctx)

	s.NoError(err)
	s.mockProducer.AssertNotCalled(s.T(), "ProduceStatsUpdate")
}

func (s *UpdatePlayerStatsSuite) TestUpdateStatsForPlayers_HandlesEmptyTrackedPlayersList() {
	s.mockStorage.On("GetAllTrackedPlayerIDs", s.ctx).Return([]int{}, nil).Once()

	err := s.playerService.UpdateStatsForPlayers(s.ctx)

	s.NoError(err)
	s.mockAPI.AssertNotCalled(s.T(), "GetPlayerByID")
	s.mockProducer.AssertNotCalled(s.T(), "ProduceStatsUpdate")
}

func (s *UpdatePlayerStatsSuite) TestUpdateStatsForPlayers_ReturnsErrorWhenGetTrackedIDsFails() {
	storageErr := errors.New("database connection lost")
	s.mockStorage.On("GetAllTrackedPlayerIDs", s.ctx).Return(nil, storageErr).Once()

	err := s.playerService.UpdateStatsForPlayers(s.ctx)

	s.Error(err)
	s.Equal(storageErr, err)
	s.mockAPI.AssertNotCalled(s.T(), "GetPlayerByID")
}

func (s *UpdatePlayerStatsSuite) TestUpdateStatsForPlayers_ContinuesOnStorageUpsertError() {
	playerIDs := []int{1, 2}

	s.mockStorage.On("GetAllTrackedPlayerIDs", s.ctx).Return(playerIDs, nil).Once()

	playerInfo1 := &models.FootballAPIPlayerInfo{ID: 1, Name: "Player1", Position: "Attacker"}
	stats1 := models.AttackerStats{GamesPlayed: 10}
	s.mockStorage.On("GetPlayerStatsByID", s.ctx, 1).Return(nil, nil).Once()
	s.mockAPI.On("GetPlayerByID", s.ctx, 1).Return(playerInfo1, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, 1, "Attacker").Return(stats1, nil).Once()
	s.mockStorage.On("UpsertPlayerStats", s.ctx, matchPlayerStats(1, "Player1")).Return(errors.New("write error")).Once()

	playerInfo2 := &models.FootballAPIPlayerInfo{ID: 2, Name: "Player2", Position: "Defender"}
	stats2 := models.DefenderStats{GamesPlayed: 10}
	s.mockStorage.On("GetPlayerStatsByID", s.ctx, 2).Return(nil, nil).Once()
	s.mockAPI.On("GetPlayerByID", s.ctx, 2).Return(playerInfo2, nil).Once()
	s.mockAPI.On("GetPlayerStatsByPosition", s.ctx, 2, "Defender").Return(stats2, nil).Once()
	s.mockStorage.On("UpsertPlayerStats", s.ctx, matchPlayerStats(2, "Player2")).Return(nil).Once()

	err := s.playerService.UpdateStatsForPlayers(s.ctx)

	s.NoError(err)
	s.mockAPI.AssertExpectations(s.T())
}

func (s *UpdatePlayerStatsSuite) TestSnapshotsEqual_ReturnsTrueForIdenticalSnapshots() {
	snap1 := &mongostorage.PlayerStats{
		ID:       1,
		Name:     "Test",
		Position: "Attacker",
		Stats:    models.AttackerStats{GamesPlayed: 10, Goals: 5},
	}
	snap2 := &mongostorage.PlayerStats{
		ID:       1,
		Name:     "Test",
		Position: "Attacker",
		Stats:    models.AttackerStats{GamesPlayed: 10, Goals: 5},
	}

	result := s.playerService.snapshotsEqual(snap1, snap2)

	s.True(result)
}

func (s *UpdatePlayerStatsSuite) TestSnapshotsEqual_ReturnsFalseForDifferentPositions() {
	snap1 := &mongostorage.PlayerStats{
		Name:     "Test",
		Position: "Attacker",
		Stats:    models.AttackerStats{GamesPlayed: 10},
	}
	snap2 := &mongostorage.PlayerStats{
		Name:     "Test",
		Position: "Defender",
		Stats:    models.AttackerStats{GamesPlayed: 10},
	}

	result := s.playerService.snapshotsEqual(snap1, snap2)

	s.True(result)
}

func (s *UpdatePlayerStatsSuite) TestSnapshotsEqual_ReturnsFalseForDifferentNames() {
	snap1 := &mongostorage.PlayerStats{
		Name:     "Test1",
		Position: "Attacker",
		Stats:    models.AttackerStats{GamesPlayed: 10},
	}
	snap2 := &mongostorage.PlayerStats{
		Name:     "Test2",
		Position: "Attacker",
		Stats:    models.AttackerStats{GamesPlayed: 10},
	}

	result := s.playerService.snapshotsEqual(snap1, snap2)

	s.False(result)
}

func (s *UpdatePlayerStatsSuite) TestSnapshotsEqual_ReturnsFalseForDifferentStats() {
	snap1 := &mongostorage.PlayerStats{
		Name:     "Test",
		Position: "Attacker",
		Stats:    models.AttackerStats{GamesPlayed: 10, Goals: 5},
	}
	snap2 := &mongostorage.PlayerStats{
		Name:     "Test",
		Position: "Attacker",
		Stats:    models.AttackerStats{GamesPlayed: 10, Goals: 6},
	}

	result := s.playerService.snapshotsEqual(snap1, snap2)

	s.False(result)
}

func (s *UpdatePlayerStatsSuite) TestSnapshotsEqual_HandlesNilSnapshots() {
	snap1 := (*mongostorage.PlayerStats)(nil)
	snap2 := (*mongostorage.PlayerStats)(nil)

	s.True(s.playerService.snapshotsEqual(snap1, snap2))

	snap3 := &mongostorage.PlayerStats{Name: "Test"}
	s.False(s.playerService.snapshotsEqual(snap1, snap3))
	s.False(s.playerService.snapshotsEqual(snap3, snap1))
}

func matchPlayerStats(expectedID int, expectedName string) *mongostorage.PlayerStats {
	return &mongostorage.PlayerStats{
		ID:   expectedID,
		Name: expectedName,
	}
}

func TestUpdatePlayerStatsSuite(t *testing.T) {
	suite.Run(t, new(UpdatePlayerStatsSuite))
}
