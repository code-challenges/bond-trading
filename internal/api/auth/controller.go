package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	. "github.com/asalvi0/bond-trading/internal/database"
	. "github.com/asalvi0/bond-trading/internal/models"
	"github.com/asalvi0/bond-trading/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Controller struct {
	db *Database
}

func newController() (*Controller, error) {
	db, err := NewDatabase()
	if err != nil {
		return nil, err
	}

	result := Controller{
		db,
	}

	return &result, nil
}

func (c *Controller) createUser(ctx context.Context, user *User) (*User, error) {
	user.CreatedAt = time.Now().UTC()

	err := c.db.CreateUser(ctx, user)
	if err != nil {
		if strings.Index(err.Error(), "SQLSTATE 23505") > -1 {
			return nil, errors.New("user already exists")
		}
		return nil, err
	}

	return user, nil
}

func (c *Controller) updateUser(ctx context.Context, user *User) error {
	updatedAt := time.Now().UTC()
	user.UpdatedAt = &updatedAt

	err := c.db.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) getUserByID(ctx context.Context, id uint) (*User, error) {
	user, err := c.db.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Controller) getUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := c.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Controller) checkPassword(ctx context.Context, email, password string) (*User, error) {
	user, err := c.getUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user.Active == false {
		return nil, errors.New("user is not active")
	}

	salt := utils.HashString(user.Email)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+salt))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
