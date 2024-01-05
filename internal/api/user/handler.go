package user

import (
	"github.com/gofiber/fiber/v2"

	"github.com/asalvi0/bond-trading/internal/api/middleware"
)

type Handler struct {
	controller Controller
}

func newHandler(controller Controller) *Handler {
	return &Handler{controller}
}

func RegisterRoutes(app *fiber.App) error {
	h := newHandler(newController())

	v1 := app.Group("/api/v1/user", middleware.Protected())

	v1.Post("/users", h.getUserOrders)
	v1.Get("/user/:id/orders", h.getUserOrders)
	v1.Get("/user/:id/orders", h.getUserOrders)

	return nil
}

func (h *Handler) getByID(c *fiber.Ctx) error {
	return c.JSON(nil)
}

func (h *Handler) getUserOrders(c *fiber.Ctx) error {
	return c.JSON(nil)
}
