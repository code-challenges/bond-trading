package order

import (
	"context"
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

func RegisterRoutes(app *fiber.App) error {
	controller, err := newController()
	if err != nil {
		return err
	}
	h := Handler{controller}

	v1 := app.Group("/api/v1/orders", middleware.Protected())

	v1.Post("/", h.createOrder)
	v1.Put("/:id", h.updateOrder)
	v1.Patch("/:id", h.cancelOrder)
	v1.Get("/", h.getOrders)
	v1.Get("/:id", h.getOrder)

	app.Get("/api/v1/my/orders", middleware.Protected(), h.getOrdersByUserID)

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

	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		return err
	}

	ctx := context.WithValue(context.Background(), "userId", userId)

	order, err := h.controller.createOrder(ctx, item)
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

	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		return err
	}

	ctx := context.WithValue(context.Background(), "userId", userId)

	order, err := h.controller.updateOrder(ctx, item)
	if err != nil {
		return err
	}

	return c.JSON(order)
}

func (h *Handler) cancelOrder(c *fiber.Ctx) error {
	id := c.Params("id", "")
	if len(id) == 0 {
		return errors.New("missing ID")
	}

	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		return err
	}

	ctx := context.WithValue(context.Background(), "userId", userId)

	err = h.controller.cancelOrder(ctx, id)
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

	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		return err
	}

	ctx := context.WithValue(context.Background(), "userId", userId)

	orders, err := h.controller.getOrders(ctx, uint(count))
	if err != nil {
		return err
	}

	return c.JSON(orders)
}

func (h *Handler) getOrder(c *fiber.Ctx) error {
	id := c.Params("id", "")
	if len(id) == 0 {
		return errors.New("missing ID")
	}

	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		return err
	}

	ctx := context.WithValue(context.Background(), "userId", userId)

	order, err := h.controller.getOrder(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(order)
}

func (h *Handler) getOrdersByUserID(c *fiber.Ctx) error {
	count := c.QueryInt("count", 0)
	if count <= 0 {
		return errors.New("invalid count provided")
	}

	userId, err := utils.GetUserIdFromToken(c)
	if err != nil {
		return err
	}

	ctx := context.WithValue(context.Background(), "userId", userId)

	orders, err := h.controller.getOrdersByUserId(ctx, userId, uint(count))
	if err != nil {
		return err
	}

	return c.JSON(orders)
}
