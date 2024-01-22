package messaging

import (
	"context"
	"fmt"

	"github.com/memphisdev/memphis.go"
)

// Create a consumer for the station(topic)
func (m *MemphisClient) createConsumer(ctx context.Context, stationName, serviceName, consumerGroupName string, callback memphis.ConsumeHandler) error {
	consumer, err := m.conn.CreateConsumer(stationName, serviceName, memphis.ConsumerGroup(consumerGroupName))
	if err != nil {
		return fmt.Errorf("Failed to create consumer: %w", err)
	}

	consumer.SetContext(ctx)
	err = consumer.Consume(callback)
	if err != nil {
		return err
	}
	m.consumers[stationName] = consumer

	return nil
}
