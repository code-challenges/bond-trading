package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Action represents the action of a bond order (Buy or Sell)
type Action string

const (
	Buy  Action = "BUY"
	Sell Action = "SELL"
)

// IsValid checks if the Action value is valid
func (a Action) IsValid() bool {
	return a == Buy || a == Sell
}

// Status represents the status of a bond order (Open, Filled, Canceled)
type Status string

const (
	Open     Status = "OPEN"
	Filled   Status = "FILLED"
	Canceled Status = "CANCELED"
)

// IsValid checks if the Status value is valid
func (s Status) IsValid() bool {
	return s == Open || s == Filled || s == Canceled
}

type Order struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id" validate:"required"`
	BondID    uint      `json:"bond_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	Price     float64   `json:"price" validate:"required,min=0.01"`
	Action    Action    `json:"action" validate:"required,oneof=BUY SELL"`
	Status    Status    `json:"status" validate:"required,oneof=OPEN FILLED CANCELED"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (order Order) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(order); err != nil {
		return err
	}
	return nil
}

func NewOrder(id, userID, bondID uint, quantity int, price float64, action Action, status Status) *Order {
	return &Order{
		ID:        id,
		UserID:    userID,
		BondID:    bondID,
		Quantity:  quantity,
		Price:     price,
		Action:    action,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // Example: Setting expiration to 24 hours from creation
	}
}
