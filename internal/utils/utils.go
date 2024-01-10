package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
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
