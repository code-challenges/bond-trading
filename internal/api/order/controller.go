package order

import (
	. "github.com/asalvi0/bond-trading/internal/database"
	"github.com/asalvi0/bond-trading/internal/messaging"
	. "github.com/asalvi0/bond-trading/internal/models"
)

type Controller struct {
	memphisClient *messaging.MemphisClient
	db            *Database
}

func newController() (*Controller, error) {
	memphisClient, err := messaging.NewMemphisClient()
	if err != nil {
		return nil, err
	}

	db, err := NewDatabase()
	if err != nil {
		return nil, err
	}

	result := Controller{
		memphisClient,
		db,
	}

	return &result, nil
}

func (c *Controller) createOrder(order *Order) (*Order, error) {
	err := c.memphisClient.ProduceMessage(order)
	if err != nil {
		return nil, err
	}

	// write to database
	go func() { c.db.CreateOrder(order) }()

	return nil, nil
}

func (c *Controller) updateOrder(order *Order) (*Order, error) {
	err := c.memphisClient.ProduceMessage(order)
	if err != nil {
		return nil, err
	}

	// write to database
	go func() { c.db.UpdateOrder(order) }()

	return nil, nil
}

func (c *Controller) cancelOrder(order *Order) error {
	err := c.memphisClient.ProduceMessage(order)
	if err != nil {
		return err
	}

	// write to database
	go func() { c.db.CancelOrder(order) }()

	return nil
}

func (c *Controller) getOrders(count uint) ([]*Order, error) {
	// read from database
	orders, err := c.db.GetOrders(count)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (c *Controller) getOrder(id string) (*Order, error) {
	// read from database
	order, err := c.db.GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (c *Controller) getOrdersByUserId(id uint, count uint) ([]*Order, error) {
	// read from database
	orders, err := c.db.GetOrdersByUserID(id, count)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
