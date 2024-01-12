package order

import (
	"context"
	"errors"
	"strings"
	"time"

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

func (c *Controller) createOrder(ctx context.Context, order *Order) (*Order, error) {
	order.Sign(ctx)

	err := c.memphisClient.ProduceMessage(order)
	if err != nil {
		return nil, err
	}

	err = c.db.CreateOrder(ctx, order)
	if err != nil {
		if strings.Index(err.Error(), "SQLSTATE 23505") > -1 {
			return nil, errors.New("order already exists")
		}
		return nil, err
	}

	return nil, nil
}

func (c *Controller) updateOrder(ctx context.Context, order *Order) (*Order, error) {
	order.UpdatedAt = time.Now().UTC()

	err := c.memphisClient.ProduceMessage(order)
	if err != nil {
		return nil, err
	}

	err = c.db.UpdateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Controller) cancelOrder(ctx context.Context, id string) error {
	order, err := c.getOrder(ctx, id)
	if err != nil {
		return err
	}

	if order.Status == CANCELLED {
		return errors.New("order already cancelled")
	}
	ogAction := order.Action
	order.Action = CANCEL // changed ONLY for the message published
	order.Status = CANCELLED
	order.UpdatedAt = time.Now().UTC()

	err = c.memphisClient.ProduceMessage(order)
	if err != nil {
		return err
	}

	order.Action = ogAction // restored to original value
	err = c.db.UpdateOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) getOrders(ctx context.Context, count uint) ([]*Order, error) {
	orders, err := c.db.GetOrders(ctx, count)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (c *Controller) getOrder(ctx context.Context, id string) (*Order, error) {
	userId := ctx.Value("userId").(uint)
	order, err := c.db.GetOrderByID(ctx, userId, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (c *Controller) getOrdersByUserId(ctx context.Context, id uint, count uint) ([]*Order, error) {
	orders, err := c.db.GetOrdersByUserID(ctx, id, count)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
