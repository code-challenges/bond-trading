package models

import (
	"time"

	"github.com/asalvi0/bond-trading/internal/utils"
	"github.com/goccy/go-json"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint       `db:"id" json:"id"`
	Email     string     `db:"email" json:"email" validate:"required,email,min=6,max=120"`
	Password  string     `db:"password_hash" json:"password" validate:"required,min=6,max=120"`
	Active    bool       `db:"active" json:"active"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
}

func (u *User) UnmarshalJSON(data []byte) error {
	type Alias User
	aux := (*Alias)(u)

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	salt := utils.HashString(u.Email)
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.CreatedAt = time.Now().UTC()
	u.Password = string(pwdHash)
	u.Active = true

	return nil
}
