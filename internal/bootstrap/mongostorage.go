package bootstrap

import (
	"log"

	"github.com/painb0w/data_processor/config"
	"github.com/painb0w/data_processor/internal/storage/mongostorage"
)

func InitMongoDb(cfg *config.Config) *mongostorage.MongoStorage {
	db, err := mongostorage.NewMongoStorage(cfg.Mongo.Uri)
	if err != nil {
		log.Panicf("error initialising MongoDB: %v", err)
	}
	return db
}
