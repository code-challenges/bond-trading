package order

import (
	"github.com/asalvi0/bond-trading/internal/messaging"
	. "github.com/asalvi0/bond-trading/internal/models"
	"github.com/asalvi0/bond-trading/internal/utils"
)

type Controller struct {
	memphisClient *messaging.MemphisClient
}

func newController() (*Controller, error) {
	memphisClient, err := messaging.NewMemphisClient()
	if err != nil {
		return nil, err
	}

	result := Controller{
		memphisClient,
	}

	return &result, nil
}

func (c *Controller) getOrdersByUserId(id uint) (orders []Order, err error) {
	return orders, err
}

func (c *Controller) createOrder(order *Order) (*Order, error) {
	// write to memphis
	msgId := utils.GenerateMessageID(order)
	err := c.memphisClient.ProduceMessage(msgId, order)
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
	err := c.memphisClient.ProduceMessage(msgId, order)
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
	err := c.memphisClient.ProduceMessage(msgId, order)
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
