package footballapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/painb0w/data_processor/internal/models"
)

func (c *FootballApiClient) getPlayerStats(ctx context.Context, id int) ([]models.FootballAPIPlayerStatistics, error) {
	url := fmt.Sprintf("%s/players?id=%d&season=2023", c.BaseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-apisports-key", c.ApiKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp models.FootballAPIPlayerResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(apiResp.Response) == 0 || len(apiResp.Response[0].Statistics) == 0 {
		return nil, fmt.Errorf("no statistics found for player %d", id)
	}

	return apiResp.Response[0].Statistics, nil
}

func (c *FootballApiClient) getGoalkeeperStats(ctx context.Context, id int) ([]models.GoalkeeperStats, error) {
	statsList, err := c.getPlayerStats(ctx, id)
	if err != nil {
		return nil, err
	}

	val := func(p *int) int {
		if p == nil {
			return 0
		}
		return *p
	}

	var result []models.GoalkeeperStats
	for _, s := range statsList {
		result = append(result, models.GoalkeeperStats{
			GamesPlayed:    val(s.Games.Appearances),
			Conceded:       val(s.Goals.Conceded),
			Saves:          val(s.Goals.Saves),
			Penalties:      val(s.Penalty.Saved) + val(s.Penalty.Missed),
			PenaltiesSaved: val(s.Penalty.Saved),
		})
	}
	return result, nil
}

func (c *FootballApiClient) getAttackerStats(ctx context.Context, id int) ([]models.AttackerStats, error) {
	statsList, err := c.getPlayerStats(ctx, id)
	if err != nil {
		return nil, err
	}

	val := func(p *int) int {
		if p == nil {
			return 0
		}
		return *p
	}

	var result []models.AttackerStats
	for _, s := range statsList {
		result = append(result, models.AttackerStats{
			GamesPlayed:     val(s.Games.Appearances),
			Shots:           val(s.Shots.Total),
			Goals:           val(s.Goals.Total),
			Assists:         val(s.Goals.Assists),
			Dribbles:        val(s.Dribbles.Attempts),
			SuccessDribbles: val(s.Dribbles.Success),
			FoulsDrawn:      val(s.Fouls.Drawn),
		})
	}
	return result, nil
}

func (c *FootballApiClient) getMidfielderStats(ctx context.Context, id int) ([]models.MidfielderStats, error) {
	statsList, err := c.getPlayerStats(ctx, id)
	if err != nil {
		return nil, err
	}

	val := func(p *int) int {
		if p == nil {
			return 0
		}
		return *p
	}

	var result []models.MidfielderStats
	for _, s := range statsList {
		result = append(result, models.MidfielderStats{
			GamesPlayed: val(s.Games.Appearances),
			Passes:      val(s.Passes.Total),
			KeyPasses:   val(s.Passes.Key),
			Assists:     val(s.Goals.Assists),
			Tackles:     val(s.Tackles.Total),
			Dribbles:    val(s.Dribbles.Attempts),
		})
	}
	return result, nil
}

func (c *FootballApiClient) getDefenderStats(ctx context.Context, id int) ([]models.DefenderStats, error) {
	statsList, err := c.getPlayerStats(ctx, id)
	if err != nil {
		return nil, err
	}

	val := func(p *int) int {
		if p == nil {
			return 0
		}
		return *p
	}

	var result []models.DefenderStats
	for _, s := range statsList {
		result = append(result, models.DefenderStats{
			GamesPlayed:   val(s.Games.Appearances),
			Tackles:       val(s.Tackles.Total),
			Interceptions: val(s.Tackles.Interceptions),
			Blocks:        val(s.Tackles.Blocks),
			FoulsCommited: val(s.Fouls.Committed),
			DuelsWon:      val(s.Duels.Won),
			Duels:         val(s.Duels.Total),
			YellowCards:   val(s.Cards.Yellow),
			RedCards:      val(s.Cards.Red) + val(s.Cards.YellowRed),
		})
	}
	return result, nil
}

func (c *FootballApiClient) GetPlayerStatsByPosition(ctx context.Context, id int, position string) (any, error) {
	pos := strings.ToLower(strings.TrimSpace(position))

	switch pos {
	case "goalkeeper":
		return c.getGoalkeeperStats(ctx, id)
	case "defender":
		return c.getDefenderStats(ctx, id)
	case "midfielder":
		return c.getMidfielderStats(ctx, id)
	case "attacker":
		return c.getAttackerStats(ctx, id)
	default:
		return nil, fmt.Errorf("unsupported position: %s", position)
	}
}
