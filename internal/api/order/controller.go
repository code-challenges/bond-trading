package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/goccy/go-json"
	"github.com/memphisdev/memphis.go"

	"github.com/asalvi0/bond-trading/internal/config"
	. "github.com/asalvi0/bond-trading/internal/models"
	"github.com/asalvi0/bond-trading/internal/utils"
)

type (
	memphisConfig struct {
		conn      *memphis.Conn
		producers map[string]*memphis.Producer
		consumers map[string]*memphis.Consumer
		stations  map[string]*memphis.Station
	}

	Controller struct {
		memphisConfig
	}
)

func newController() (*Controller, error) {
	result := Controller{}

	result.producers = make(map[string]*memphis.Producer, 3)
	result.consumers = make(map[string]*memphis.Consumer, 3)
	result.stations = make(map[string]*memphis.Station, 3)

	err := result.setupMemphis()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Controller) getOrdersByUserId(id uint) (orders []Order, err error) {
	return orders, err
}

func (c *Controller) createOrder(order *Order) (*Order, error) {
	// write to memphis
	msgId := utils.GenerateMessageID(order)
	err := c.produceMessage(msgId, order)
	if err != nil {
		return nil, err
	}

	// write to database
	go func() {
		// err := order.Insert()
		// if err != nil {
		// 	return nil, err
		// }
	}()

	return nil, nil
}

func (c *Controller) updateOrder(order *Order) (*Order, error) {
	// write to memphis
	msgId := utils.GenerateMessageID(order)
	err := c.produceMessage(msgId, order)
	if err != nil {
		return nil, err
	}

	// write to database
	go func() {
		// err := order.Insert()
		// if err != nil {
		// 	return nil, err
		// }
	}()

	return nil, nil
}

func (c *Controller) cancelOrder(order *Order) error {
	// write to memphis
	msgId := utils.GenerateMessageID(order)
	err := c.produceMessage(msgId, order)
	if err != nil {
		return err
	}

	// write to database
	go func() {
		// err := order.Insert()
		// if err != nil {
		// 	return nil, err
		// }
	}()

	return nil
}

func (c *Controller) getOrders(count uint) (orders []Order, err error) {
	return orders, err
}

func (c *Controller) getOrder(id uint) (order *Order, err error) {
	return order, err
}

// Initialize memphis connection, stations(topics), producers and consumers
func (c *Controller) setupMemphis() (err error) {
	host := config.Config("MEMPHIS_HOST")
	user := config.Config("MEMPHIS_USERNAME")
	password := config.Config("MEMPHIS_PASSWORD")

	c.conn, err = memphis.Connect(host, user, memphis.Password(password))
	if err != nil {
		return err
	}

	c.setupStations()
	c.setupProducers()
	c.setupConsumers()

	return nil
}

// Create the buy and sell stations(topics)
func (c *Controller) setupStations() error {
	buyStationName := config.Config("MEMPHIS_BUY_STATION_NAME")
	sellStationName := config.Config("MEMPHIS_SELL_STATION_NAME")
	cancelStationName := config.Config("MEMPHIS_CANCEL_STATION_NAME")

	// CANCEL
	err := c.createStation(cancelStationName)
	if err != nil {
		return err
	}

	// SELL
	err = c.createStation(sellStationName)
	if err != nil {
		return err
	}

	// BUY
	err = c.createStation(buyStationName)
	if err != nil {
		return err
	}

	return nil
}

// Create the buy, sell and cancel producers
func (c *Controller) setupProducers() error {
	serviceName := config.Config("MEMPHIS_PRODUCER_SERVICE_NAME")
	buyStationName := config.Config("MEMPHIS_BUY_STATION_NAME")
	sellStationName := config.Config("MEMPHIS_SELL_STATION_NAME")
	cancelStationName := config.Config("MEMPHIS_CANCEL_STATION_NAME")

	// BUY
	err := c.createProducer(serviceName, buyStationName)
	if err != nil {
		return fmt.Errorf("Failed to create producer: %w", err)
	}

	// SELL
	err = c.createProducer(serviceName, sellStationName)
	if err != nil {
		return fmt.Errorf("Failed to create producer: %w", err)
	}

	// CANCEL
	err = c.createProducer(serviceName, cancelStationName)
	if err != nil {
		return fmt.Errorf("Failed to create producer: %w", err)
	}

	return nil
}

// Create the buy, sell and cancel consumers
func (c *Controller) setupConsumers() error {
	serviceName := config.Config("MEMPHIS_CONSUMER_SERVICE_NAME")

	buyStationName := config.Config("MEMPHIS_BUY_STATION_NAME")
	sellStationName := config.Config("MEMPHIS_SELL_STATION_NAME")
	cancelStationName := config.Config("MEMPHIS_CANCEL_STATION_NAME")

	buyConsumerGroupName := config.Config("MEMPHIS_BUY_CONSUMER_GROUP_NAME")
	sellConsumerGroupName := config.Config("MEMPHIS_SELL_CONSUMER_GROUP_NAME")
	cancelConsumerGroupName := config.Config("MEMPHIS_CANCEL_CONSUMER_GROUP_NAME")

	ctx := context.Background()

	// BUY
	err := c.createConsumer(ctx, buyStationName, serviceName, buyConsumerGroupName)
	if err != nil {
		return fmt.Errorf("Failed to create consumer: %w", err)
	}

	// SELL
	err = c.createConsumer(ctx, sellStationName, serviceName, sellConsumerGroupName)
	if err != nil {
		return fmt.Errorf("Failed to create consumer: %w", err)
	}

	// CANCEL
	err = c.createConsumer(ctx, cancelStationName, serviceName, cancelConsumerGroupName)
	if err != nil {
		return fmt.Errorf("Failed to create consumer: %w", err)
	}

	return nil
}

// Create a memphis station(topic) and it's respective dead letter station for re-consumption
func (c *Controller) createStation(name string) error {
	opts := []memphis.StationOpt{
		memphis.StorageTypeOpt(memphis.Disk),
		memphis.RetentionTypeOpt(memphis.MaxMessageAgeSeconds),
		memphis.RetentionVal(int((24 * time.Hour).Seconds())),
	}

	dlStationName := name + "_dead_letter"
	dlStation, err := c.conn.CreateStation(dlStationName, opts...)
	if err != nil {
		return err
	}
	c.stations[name] = dlStation

	opts = append(opts, memphis.DlsStation(dlStationName))
	station, err := c.conn.CreateStation(name, opts...)
	if err != nil {
		return err
	}
	c.stations[name] = station

	return nil
}

// Create a producer for the station(topic)
func (c *Controller) createProducer(serviceName, stationName string) error {
	producer, err := c.conn.CreateProducer(stationName, serviceName)
	if err != nil {
		return fmt.Errorf("Failed to create producer: %w", err)
	}
	c.producers[stationName] = producer

	return nil
}

// Create a consumer for the station(topic)
func (c *Controller) createConsumer(ctx context.Context, stationName, serviceName, consumerGroupName string) error {
	consumer, err := c.conn.CreateConsumer(stationName, serviceName, memphis.ConsumerGroup(consumerGroupName))
	if err != nil {
		return fmt.Errorf("Failed to create consumer: %w", err)
	}

	consumer.SetContext(ctx)
	err = consumer.Consume(c.consumeMessages)
	if err != nil {
		return err
	}
	c.consumers[stationName] = consumer

	return nil
}

// Send a message to the memphis station(topic) depending on the action provided
func (c *Controller) produceMessage(msgId string, order *Order) error {
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

	producer, ok := c.producers[station]
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
func (c *Controller) consumeMessages(msgs []*memphis.Msg, err error, ctx context.Context) {
	if err != nil {
		log.Printf("Fetch failed: %v", err)
		return
	}

	for i := 0; i < len(msgs); i++ {
		log.Println(string(msgs[i].Data()))
		msgs[i].Ack()
	}
}
