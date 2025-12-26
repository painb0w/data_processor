package mongostorage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoStorage) UpsertPlayerStats(ctx context.Context, snapshot *PlayerStats) error {
	snapshot.UpdatedAt = time.Now()

	collection := m.client.Database(dbName).Collection(collectionName)

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": snapshot.ID},
		bson.M{"$set": snapshot},
		options.Update().SetUpsert(true),
	)
	return err
}
