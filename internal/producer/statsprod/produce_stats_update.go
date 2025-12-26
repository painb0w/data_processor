package statsprod

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type StatsUpdateEvent struct {
	PlayerID int64 `json:"player_id"`
}

func (p *StatsProducer) ProduceStatsUpdate(ctx context.Context, playerID int) error {
	event := StatsUpdateEvent{
		PlayerID: int64(playerID),
	}

	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal stats event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", playerID)),
		Value: value,
		Time:  time.Now(),
	}

	return p.writer.WriteMessages(ctx, msg)
}
