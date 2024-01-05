package models

import "time"

// Action represents the action of a bond order (Buy, Sell or Cancel)
type Action string

const (
	Buy    Action = "BUY"
	Sell   Action = "SELL"
	Cancel Action = "CANCEL"
)

// IsValid checks if the Action value is valid
func (a Action) IsValid() bool {
	return a == Buy || a == Sell || a == Cancel
}

// Status represents the status of a bond order (Open, Filled, Canceled)
type Status string

const (
	Pending  Status = "PENDING"
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
	UserID    uint      `json:"userId" validate:"required"`
	BondID    uint      `json:"bondId" validate:"required"`
	Quantity  uint      `json:"quantity" validate:"required,min=1,max=10000"`
	Price     float64   `json:"price" validate:"required,min=0,max=100000000"`
	Action    Action    `json:"action" validate:"required,oneof=BUY SELL CANCEL"`
	Status    Status    `json:"status" validate:"required,oneof=PENDING OPEN FILLED CANCELED"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewOrder(id, userID, bondID, quantity uint, price float64, action Action, status Status) *Order {
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
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
}
