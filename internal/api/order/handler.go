package order

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/asalvi0/bond-trading/internal/api/middleware"
	"github.com/asalvi0/bond-trading/internal/models"
	"github.com/asalvi0/bond-trading/internal/utils"
)

type Handler struct {
	controller *Controller
}

func newHandler(controller *Controller) *Handler {
	return &Handler{controller}
}

func RegisterRoutes(app *fiber.App) error {
	controller, err := newController()
	if err != nil {
		return err
	}
	h := newHandler(controller)

	v1 := app.Group("/api/v1/orders", middleware.Protected())

	v1.Post("/", h.createOrder)
	v1.Put("/", h.updateOrder)
	v1.Put("/:id", h.cancelOrder)
	v1.Get("/", h.getOrders)
	v1.Get("/:id", h.getOrder)

	return nil
}

func (h *Handler) createOrder(c *fiber.Ctx) error {
	item := new(models.Order)
	if err := c.BodyParser(item); err != nil {
		return err
	}

	err := utils.ValidateInput(item)
	if err != nil {
		return err
	}

	order, err := h.controller.createOrder(item)
	if err != nil {
		return err
	}

	return c.JSON(order)
}

func (h *Handler) updateOrder(c *fiber.Ctx) error {
	item := new(models.Order)
	if err := c.BodyParser(item); err != nil {
		return err
	}

	err := utils.ValidateInput(item)
	if err != nil {
		return err
	}

	order, err := h.controller.updateOrder(item)
	if err != nil {
		return err
	}

	return c.JSON(order)
}

func (h *Handler) cancelOrder(c *fiber.Ctx) error {
	order := new(models.Order)
	if err := c.BodyParser(order); err != nil {
		return err
	}

	err := utils.ValidateInput(order)
	if err != nil {
		return err
	}

	err = h.controller.cancelOrder(order)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) getOrders(c *fiber.Ctx) error {
	count := c.QueryInt("count", 0)
	if count <= 0 {
		return errors.New("invalid count provided")
	}

	orders, err := h.controller.getOrders(uint(count))
	if err != nil {
		return err
	}

	return c.JSON(orders)
}

func (h *Handler) getOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	} else if id <= 0 {
		return errors.New("invalid ID provided")
	}

	order, err := h.controller.getOrder(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(order)
}
