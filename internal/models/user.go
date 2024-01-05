package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username" validate:"required,min=3,max=50"`
	Email        string    `json:"email" validate:"required,email"`
	PasswordHash string    `json:"password" validate:"required,min=8"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func NewUser(id uint, username, email, password string) (*User, error) {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: string(pwdHash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
