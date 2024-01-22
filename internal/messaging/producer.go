package messaging

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/memphisdev/memphis.go"

	"github.com/asalvi0/bond-trading/internal/config"
	. "github.com/asalvi0/bond-trading/internal/models"
)

// Create a producer for the station(topic)
func (m *MemphisClient) createProducer(serviceName, stationName string) error {
	producer, err := m.conn.CreateProducer(stationName, serviceName)
	if err != nil {
		return fmt.Errorf("Failed to create producer: %w", err)
	}
	m.producers[stationName] = producer

	return nil
}

// Send a message to the memphis station(topic) depending on the action provided
func (m *MemphisClient) ProduceMessage(order *Order) error {
	var station string
	switch order.Action {
	case BUY:
		station = config.Config("MEMPHIS_BUY_STATION_NAME")
	case SELL:
		station = config.Config("MEMPHIS_SELL_STATION_NAME")
	case CANCEL:
		station = config.Config("MEMPHIS_CANCEL_STATION_NAME")
	default:
		return errors.New("Invalid action provided")
	}

	producer, ok := m.producers[station]
	if producer == nil || !ok {
		return errors.New("Producer not found")
	}

	opts := []memphis.ProduceOpt{
		memphis.AsyncProduce(),
		memphis.AckWaitSec(15),
		memphis.MsgId(order.ID),
	}

	jsonData, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = producer.Produce(jsonData, opts...)
	if err != nil {
		return fmt.Errorf("Produce failed: %w", err)
	}

	return nil
}
