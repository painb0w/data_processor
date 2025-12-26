package mongostorage

import "time"

type PlayerStats struct {
	ID        int       `bson:"_id"`
	Name      string    `bson:"name"`
	Position  string    `bson:"position"`
	Stats     any       `bson:"stats"`
	UpdatedAt time.Time `bson:"updated_at"`
}

const (
	dbName         = "football_stats"
	collectionName = "players_snapshots"
)
