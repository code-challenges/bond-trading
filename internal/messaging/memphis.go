package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/memphisdev/memphis.go"

	"github.com/asalvi0/bond-trading/internal/config"
)

type MemphisClient struct {
	conn      *memphis.Conn
	stations  map[string]*memphis.Station
	producers map[string]*memphis.Producer
	consumers map[string]*memphis.Consumer
}

func NewMemphisClient() (*MemphisClient, error) {
	result := MemphisClient{}
	result.producers = make(map[string]*memphis.Producer, 3)
	result.consumers = make(map[string]*memphis.Consumer, 3)
	result.stations = make(map[string]*memphis.Station, 3)

	err := result.setupMemphis()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *MemphisClient) Close() {
	m.conn.Close()
}

// Initialize memphis connection, stations(topics), producers and consumers
func (m *MemphisClient) setupMemphis() (err error) {
	host := config.Config("MEMPHIS_HOST")
	user := config.Config("MEMPHIS_USERNAME")
	password := config.Config("MEMPHIS_PASSWORD")

	m.conn, err = memphis.Connect(host, user, memphis.Password(password))
	if err != nil {
		return err
	}

	m.setupStations()
	// m.setupProducers()
	// m.setupConsumers()

	return nil
}

// Create the buy and sell stations(topics)
func (m *MemphisClient) setupStations() error {
	buyStationName := config.Config("MEMPHIS_BUY_STATION_NAME")
	sellStationName := config.Config("MEMPHIS_SELL_STATION_NAME")
	cancelStationName := config.Config("MEMPHIS_CANCEL_STATION_NAME")

	// CANCEL
	err := m.createStation(cancelStationName)
	if err != nil {
		return err
	}

	// SELL
	err = m.createStation(sellStationName)
	if err != nil {
		return err
	}

	// BUY
	err = m.createStation(buyStationName)
	if err != nil {
		return err
	}

	return nil
}

// Create the buy, sell and cancel producers
func (m *MemphisClient) SetupProducers() error {
	serviceName := config.Config("MEMPHIS_PRODUCER_SERVICE_NAME")

	buyStationName := config.Config("MEMPHIS_BUY_STATION_NAME")
	sellStationName := config.Config("MEMPHIS_SELL_STATION_NAME")
	cancelStationName := config.Config("MEMPHIS_CANCEL_STATION_NAME")

	// BUY
	err := m.createProducer(serviceName, buyStationName)
	if err != nil {
		return fmt.Errorf("Failed to create producer: %w", err)
	}

	// SELL
	err = m.createProducer(serviceName, sellStationName)
	if err != nil {
		return fmt.Errorf("Failed to create producer: %w", err)
	}

	// CANCEL
	err = m.createProducer(serviceName, cancelStationName)
	if err != nil {
		return fmt.Errorf("Failed to create producer: %w", err)
	}

	return nil
}

// Create the buy, sell and cancel consumers
func (m *MemphisClient) SetupConsumers(callback memphis.ConsumeHandler) error {
	serviceName := config.Config("MEMPHIS_CONSUMER_SERVICE_NAME")

	buyStationName := config.Config("MEMPHIS_BUY_STATION_NAME")
	sellStationName := config.Config("MEMPHIS_SELL_STATION_NAME")
	cancelStationName := config.Config("MEMPHIS_CANCEL_STATION_NAME")

	buyConsumerGroupName := config.Config("MEMPHIS_BUY_CONSUMER_GROUP_NAME")
	sellConsumerGroupName := config.Config("MEMPHIS_SELL_CONSUMER_GROUP_NAME")
	cancelConsumerGroupName := config.Config("MEMPHIS_CANCEL_CONSUMER_GROUP_NAME")

	ctx := context.Background()

	// BUY
	err := m.createConsumer(ctx, buyStationName, serviceName, buyConsumerGroupName, callback)
	if err != nil {
		return fmt.Errorf("Failed to create consumer: %w", err)
	}

	// SELL
	err = m.createConsumer(ctx, sellStationName, serviceName, sellConsumerGroupName, callback)
	if err != nil {
		return fmt.Errorf("Failed to create consumer: %w", err)
	}

	// CANCEL
	err = m.createConsumer(ctx, cancelStationName, serviceName, cancelConsumerGroupName, callback)
	if err != nil {
		return fmt.Errorf("Failed to create consumer: %w", err)
	}

	return nil
}

// Create a memphis station(topic) and it's respective dead letter station for re-consumption
func (m *MemphisClient) createStation(name string) error {
	opts := []memphis.StationOpt{
		memphis.StorageTypeOpt(memphis.Disk),
		memphis.RetentionTypeOpt(memphis.MaxMessageAgeSeconds),
		memphis.RetentionVal(int((10 * time.Hour).Seconds())),
	}

	dlStationName := name + "_dead_letter"
	dlStation, err := m.conn.CreateStation(dlStationName, opts...)
	if err != nil {
		return err
	}
	m.stations[name] = dlStation

	opts = append(opts, memphis.DlsStation(dlStationName))
	station, err := m.conn.CreateStation(name, opts...)
	if err != nil {
		return err
	}
	m.stations[name] = station

	return nil
}
