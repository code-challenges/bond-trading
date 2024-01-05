package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/memphisdev/memphis.go"

	"github.com/asalvi0/bond-trading/internal/config"
	. "github.com/asalvi0/bond-trading/internal/models"
)

type MemphisClient struct {
	conn      *memphis.Conn
	producers map[string]*memphis.Producer
	consumers map[string]*memphis.Consumer
	stations  map[string]*memphis.Station
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
	m.setupProducers()
	m.setupConsumers()

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
func (m *MemphisClient) setupProducers() error {
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
func (m *MemphisClient) setupConsumers() error {
	serviceName := config.Config("MEMPHIS_CONSUMER_SERVICE_NAME")

	buyStationName := config.Config("MEMPHIS_BUY_STATION_NAME")
	sellStationName := config.Config("MEMPHIS_SELL_STATION_NAME")
	cancelStationName := config.Config("MEMPHIS_CANCEL_STATION_NAME")

	buyConsumerGroupName := config.Config("MEMPHIS_BUY_CONSUMER_GROUP_NAME")
	sellConsumerGroupName := config.Config("MEMPHIS_SELL_CONSUMER_GROUP_NAME")
	cancelConsumerGroupName := config.Config("MEMPHIS_CANCEL_CONSUMER_GROUP_NAME")

	ctx := context.Background()

	// BUY
	err := m.createConsumer(ctx, buyStationName, serviceName, buyConsumerGroupName)
	if err != nil {
		return fmt.Errorf("Failed to create consumer: %w", err)
	}

	// SELL
	err = m.createConsumer(ctx, sellStationName, serviceName, sellConsumerGroupName)
	if err != nil {
		return fmt.Errorf("Failed to create consumer: %w", err)
	}

	// CANCEL
	err = m.createConsumer(ctx, cancelStationName, serviceName, cancelConsumerGroupName)
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
		memphis.RetentionVal(int((24 * time.Hour).Seconds())),
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

// Create a producer for the station(topic)
func (m *MemphisClient) createProducer(serviceName, stationName string) error {
	producer, err := m.conn.CreateProducer(stationName, serviceName)
	if err != nil {
		return fmt.Errorf("Failed to create producer: %w", err)
	}
	m.producers[stationName] = producer

	return nil
}

// Create a consumer for the station(topic)
func (m *MemphisClient) createConsumer(ctx context.Context, stationName, serviceName, consumerGroupName string) error {
	consumer, err := m.conn.CreateConsumer(stationName, serviceName, memphis.ConsumerGroup(consumerGroupName))
	if err != nil {
		return fmt.Errorf("Failed to create consumer: %w", err)
	}

	consumer.SetContext(ctx)
	err = consumer.Consume(m.ConsumeMessages)
	if err != nil {
		return err
	}
	m.consumers[stationName] = consumer

	return nil
}

// Send a message to the memphis station(topic) depending on the action provided
func (m *MemphisClient) ProduceMessage(msgId string, order *Order) error {
	var station string
	switch order.Action {
	case Buy:
		station = config.Config("MEMPHIS_BUY_STATION_NAME")
	case Sell:
		station = config.Config("MEMPHIS_SELL_STATION_NAME")
	case Cancel:
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
		memphis.MsgId(msgId),
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

// Process a batch of messages from a memphis station(topic)
func (m *MemphisClient) ConsumeMessages(msgs []*memphis.Msg, err error, ctx context.Context) {
	if err != nil {
		log.Printf("Fetch failed: %v", err)
		return
	}

	for i := 0; i < len(msgs); i++ {
		log.Println(string(msgs[i].Data()))
		msgs[i].Ack()
	}
}
