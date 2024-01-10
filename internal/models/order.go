package models

import (
	"time"

	"github.com/asalvi0/bond-trading/internal/utils"
)

// Action represents the action of a bond order (Buy, Sell or Cancel)
type Action string

const (
	BUY    Action = "BUY"
	SELL   Action = "SELL"
	CANCEL Action = "CANCEL"
)

// IsValid checks if the Action value is valid
func (a Action) IsValid() bool {
	return a == BUY || a == SELL || a == CANCEL
}

func (a Action) ToSide() int {
	switch a {
	case SELL:
		return 0
	case BUY:
		return 1
	default:
		return -1
	}
}

// Status represents the status of a bond order (Open, Filled, Canceled)
type Status string

const (
	PENDING  Status = "PENDING"
	OPEN     Status = "OPEN"
	FILLED   Status = "FILLED"
	CANCELED Status = "CANCELED"
)

// IsValid checks if the Status value is valid
func (s Status) IsValid() bool {
	return s == OPEN || s == FILLED || s == CANCELED || s == PENDING
}

type Order struct {
	ID        string    `json:"id"`
	UserID    uint      `json:"userId"`
	BondID    uint      `json:"bondId" validate:"required"`
	Quantity  uint      `json:"quantity" validate:"required,min=1,max=10000"`
	Filled    uint      `json:"filled" validate:"min=0,max=10000"`
	Price     float32   `json:"price" validate:"required,min=0,max=100000000"`
	Action    Action    `json:"action" validate:"required,oneof=BUY SELL CANCEL"`
	Status    Status    `json:"status" validate:"required,oneof=PENDING OPEN FILLED CANCELED"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewOrder(bondID, quantity uint, price float32, action Action) *Order {
	order := Order{
		BondID:    bondID,
		Quantity:  quantity,
		Price:     price,
		Action:    action,
		Status:    PENDING,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(10 * time.Hour),
	}
	order.ID = utils.GenerateID(order)

	return &order
}
