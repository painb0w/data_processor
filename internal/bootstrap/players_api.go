package bootstrap

import (
	server "github.com/painb0w/data_processor/internal/api/player_service_api"
	"github.com/painb0w/data_processor/internal/services/playersservice"
)

func InitPlayerServiceAPI(ps *playersservice.PlayerService) *server.PlayersServiceAPI {
	return server.NewPlayersServiceAPI(ps)
}
