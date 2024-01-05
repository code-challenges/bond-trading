package user

import model "github.com/asalvi0/bond-trading/internal/models"

type Controller struct{}

func newController() Controller {
	return Controller{}
}

func (c *Controller) getByUserId(id uint) ([]model.Order, error) {
	return nil, nil
}

func (c *Controller) createUser(order *model.Order) error {
	return nil
}

func (c *Controller) updateUser(order *model.Order) error {
	return nil
}

func (c *Controller) deleteUserById(order *model.Order) error {
	return nil
}
