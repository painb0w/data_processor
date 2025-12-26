package mongostorage

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoStorage) GetPlayerStatsByID(ctx context.Context, playerID int) (*PlayerStats, error) {
	collection := m.client.Database(dbName).Collection(collectionName)

	var snapshot PlayerStats
	err := collection.FindOne(ctx, bson.M{"_id": playerID}).Decode(&snapshot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("player not found")
		}
		return nil, err
	}

	return &snapshot, nil
}

func (m *MongoStorage) GetAllTrackedPlayerIDs(ctx context.Context) ([]int, error) {
	collection := m.client.Database(dbName).Collection(collectionName)
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ids []int
	for cursor.Next(ctx) {
		var item struct {
			ID int `bson:"_id"`
		}
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		ids = append(ids, item.ID)
	}
	return ids, nil
}
