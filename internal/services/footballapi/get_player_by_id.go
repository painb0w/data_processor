package footballapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/painb0w/data_processor/internal/models"
)

func (c *FootballApiClient) GetPlayerByID(ctx context.Context, id int) (*models.FootballAPIPlayerInfo, error) {
	url := fmt.Sprintf("%s/players/profiles?player=%d", c.BaseURL, id)
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

	var apiResponse struct {
		Response []struct {
			Player models.FootballAPIPlayerInfo `json:"player"`
		} `json:"response"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse.Response[0].Player, nil
}
