package bootstrap

import (
	"fmt"
	"log"
	"strconv"

	"github.com/painb0w/data_processor/config"
	"github.com/painb0w/data_processor/internal/storage/redisstorage"
)

func InitRedisStorage(cfg *config.Config) *redisstorage.RedisStorage {

	db, err := strconv.Atoi(cfg.Redis.Db)
	if err != nil {
		db = 0
	}

	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)

	storage, err := redisstorage.NewRedisStorage(addr, cfg.Redis.Password, db)
	if err != nil {
		log.Panicf("error initialising Redis, %v", err)
		panic(err)
	}
	return storage
}
