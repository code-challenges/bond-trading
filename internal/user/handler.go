package user

import (
	"github.com/asalvi0/bond-trading/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1/user", middleware.Protected())

	v1.Post("/users", getUserOrders)
	v1.Get("/user/:id/orders", getUserOrders)
	v1.Get("/user/:id/orders", getUserOrders)
}

func getByID(c *fiber.Ctx) error {
	return c.JSON(nil)
}

func getUserOrders(c *fiber.Ctx) error {
	return c.JSON(nil)
}
