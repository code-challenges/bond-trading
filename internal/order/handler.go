package order

import (
	"github.com/asalvi0/bond-trading/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1/order", middleware.Protected())

	v1.Get("/orders", getAllOrders)
	v1.Get("/order/:id", getOrder)
	v1.Post("/order", createOrder)
	v1.Put("/order/:id", updateOrder)
	v1.Delete("/order/:id", cancelOrder)
}

func createOrder(c *fiber.Ctx) error {
	return c.JSON(nil)
}

func updateOrder(c *fiber.Ctx) error {
	return c.JSON(nil)
}

func cancelOrder(c *fiber.Ctx) error {
	return c.JSON(nil)
}

func getOrder(c *fiber.Ctx) error {
	return c.JSON(nil)
}

func getAllOrders(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}
