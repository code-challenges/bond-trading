package auth

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/asalvi0/bond-trading/internal/config"
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

	v1 := app.Group("/api/v1/auth")

	v1.Post("/signin", h.login)
	v1.Post("/signup", h.signup)
	v1.Patch("/reset-password", h.resetPassword)

	return nil
}

func (h *Handler) signup(c *fiber.Ctx) error {
	input := new(models.User)
	if err := c.BodyParser(input); err != nil {
		return err
	}

	err := utils.ValidateInput(input)
	if err != nil {
		return err
	}

	ctx := context.Background()
	user, err := h.controller.createUser(ctx, input)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (h *Handler) resetPassword(c *fiber.Ctx) error {
	input := new(models.User)
	if err := c.BodyParser(input); err != nil {
		return err
	}

	err := utils.ValidateInput(input)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = h.controller.updateUser(ctx, input)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *Handler) login(c *fiber.Ctx) error {
	type Input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	input := new(Input)
	if err := c.BodyParser(input); err != nil {
		return err
	}

	_, err := mail.ParseAddress(input.Email)
	if err != nil {
		return errors.New("invalid email")
	}

	ctx := context.Background()
	user, err := h.controller.checkPassword(ctx, input.Email, input.Password)
	if err != nil {
		return err
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	signedAccessToken, err := accessToken.SignedString([]byte(config.Config("JWT_SECRET_KEY")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claims = refreshToken.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID

	signedRefreshToken, err := refreshToken.SignedString([]byte(config.Config("JWT_SECRET_KEY")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"accessToken":  signedAccessToken,
		"refreshToken": signedRefreshToken,
	})
}
