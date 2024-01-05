package authn

import (
	"github.com/asalvi0/bond-trading/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) error {
	v1 := app.Group("/v1")

	// GET
	v1.Get("/email/verify/:id/:hash", VerifyEmail)
	v1.Get("/resend/:id", Resend)

	// POST
	v1.Post("/logout", Logout, middleware.Protected())
	v1.Post("/login", Login)
	v1.Post("/register", Register)
	v1.Post("/reset-confirm", ResetConfirm)
	v1.Post("/reset-password", ResetPassword)
	v1.Post("/user/change-password", PasswordChange, middleware.Protected())
	v1.Post("/email/resend", EmailResend)

	return nil
}

// Logout - Cierre de sesión
func Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

// Login - Inicio de sesión
func Login(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

// Register - Registro de un cliente
func Register(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

// ResetConfirm - Segundo paso de la recuperación de contraseña.
func ResetConfirm(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

// ResetPassword - Envía un correo al usuario para que pueda recuperar su contraseña.
func ResetPassword(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

// PasswordChange - Actualizar información de un usuario
func PasswordChange(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

// EmailResend - Envía un correo de verificación
func EmailResend(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

// VerifyEmail - Verifica un usuario
func VerifyEmail(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

// Resend - Obtener URL para reenviar el correo
func Resend(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
