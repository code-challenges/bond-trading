package models

import (
	"context"
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
	PENDING   Status = "PENDING"
	OPEN      Status = "OPEN"
	FILLED    Status = "FILLED"
	CANCELLED Status = "CANCELLED"
)

// IsValid checks if the Status value is valid
func (s Status) IsValid() bool {
	return s == OPEN || s == FILLED || s == CANCELLED || s == PENDING
}

type Order struct {
	ID        string    `db:"id" json:"id"`
	UserID    uint      `db:"user_id" json:"userId"`
	BondID    uint      `db:"bond_id" json:"bondId" validate:"required"`
	Quantity  uint      `db:"quantity" json:"quantity" validate:"required,min=1,max=10000"`
	Filled    uint      `db:"filled" json:"filled" validate:"min=0,max=10000"`
	Price     float32   `db:"price" json:"price" validate:"required,min=0,max=100000000"`
	Action    Action    `db:"action" json:"action" validate:"required,oneof=BUY SELL CANCEL"`
	Status    Status    `db:"status" json:"status" validate:"required,oneof=PENDING OPEN FILLED CANCELLED"`
	ExpiresAt time.Time `db:"expires_at" json:"expiresAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

func (o *Order) Sign(ctx context.Context) {
	o.UserID = ctx.Value("userId").(uint)
	o.CreatedAt = time.Now().UTC()
	o.ExpiresAt = o.CreatedAt.Add(10 * time.Hour)
	o.ID = utils.GenerateID(o) // TODO:  check for duplicates and move it at the top to enable idempotency
}

func NewOrder(bondID, quantity uint, price float32, action Action) *Order {
	order := Order{
		BondID:    bondID,
		Quantity:  quantity,
		Price:     price,
		Action:    action,
		Status:    PENDING,
		ExpiresAt: time.Now().UTC().Add(10 * time.Hour),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	order.ID = utils.GenerateID(order)

	return &order
}
