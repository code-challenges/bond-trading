package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/asalvi0/bond-trading/internal/utils"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username" validate:"required,min=3,max=50"`
	Email        string    `json:"email" validate:"required,email"`
	PasswordHash string    `json:"password" validate:"required,min=8"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func NewUser(username, email, password string) (*User, error) {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := User{
		Username:     username,
		Email:        email,
		PasswordHash: string(pwdHash),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	user.ID = utils.GenerateID(user)

	return &user, nil
}
