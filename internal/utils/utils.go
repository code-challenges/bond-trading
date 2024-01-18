package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zeebo/xxh3"
)

func ValidateInput[T any](item T) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(item); err != nil {
		return err
	}
	return nil
}

func GenerateID[T any](data T) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	hash := xxh3.Hash(jsonData)
	id := fmt.Sprintf("%x", hash)

	return id
}

func HashString(data string) string {
	hash := xxh3.Hash([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func GetUserIdFromToken(c *fiber.Ctx) (uint, error) {
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return 0, errors.New("missing user token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid user token")
	}

	uidClaim, ok := claims["uid"]
	if !ok {
		return 0, errors.New("invalid user token")
	}

	uidFloat, ok := uidClaim.(float64)
	if !ok {
		return 0, errors.New("invalid user token")
	}
	userId := uint(uidFloat)

	return userId, nil
}
