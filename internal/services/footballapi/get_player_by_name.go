package footballapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/painb0w/data_processor/internal/models"
)

func (c *FootballApiClient) GetPlayerByName(ctx context.Context, name string) ([]models.FootballAPIPlayerInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/players/profiles?search="+name, nil)
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

	var apiResponse struct {
		Response []struct {
			Player models.FootballAPIPlayerInfo `json:"player"`
		} `json:"response"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, err
	}

	players := make([]models.FootballAPIPlayerInfo, len(apiResponse.Response))
	for i, item := range apiResponse.Response {
		players[i] = item.Player
	}

	return players, nil
}
