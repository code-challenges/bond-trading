package database

import (
	"context"
	"errors"
	"fmt"

	. "github.com/asalvi0/bond-trading/internal/models"
)

func (db *Database) CreateOrder(ctx context.Context, order *Order) error {
	sql := `INSERT INTO orders (id, user_id, bond_id, quantity, filled, price, action, status, expires_at, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	tag, err := db.dbPool.Exec(ctx, sql,
		order.ID, order.UserID, order.BondID, order.Quantity, order.Filled,
		order.Price, order.Action, order.Status, order.ExpiresAt, order.CreatedAt)
	if err != nil {
		return fmt.Errorf("Unable to insert record: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("No records were created")
	}

	return nil
}

func (db *Database) UpdateOrder(ctx context.Context, order *Order) error {
	sql := `UPDATE orders SET
				quantity = $4,
				filled = $5,
				price = $6,
				expires_at = $7,
				updated_at = $8
			WHERE id = $1 AND user_id = $2 AND bond_id = $3`

	tag, err := db.dbPool.Exec(ctx, sql,
		order.ID, order.UserID, order.BondID, order.Quantity, order.Filled,
		order.Price, order.ExpiresAt, order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("Unable to update record: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("No records were updated")
	}

	return nil
}

func (db *Database) UpdateOrderStatus(ctx context.Context, order *Order) error {
	sql := `UPDATE orders SET status = $4, updated_at = $5
			WHERE id = $1 AND user_id = $2 AND bond_id = $3`

	tag, err := db.dbPool.Exec(ctx, sql,
		order.ID, order.UserID, order.BondID, order.Status, order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("Unable to cancel order: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return errors.New("No records were updated")
	}

	return nil
}

func (db *Database) GetOrders(ctx context.Context, count uint) ([]*Order, error) {
	sql := `SELECT * FROM orders ORDER BY created_at DESC LIMIT $1`

	rows, err := db.dbPool.Query(ctx, sql, count)
	if err != nil {
		return nil, fmt.Errorf("Unable to select records: %w", err)
	}
	defer rows.Close()

	orders := make([]*Order, 0)
	for rows.Next() {
		var order Order
		err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.BondID,
			&order.Quantity,
			&order.Filled,
			&order.Price,
			&order.Action,
			&order.Status,
			&order.ExpiresAt,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			break
		}
		orders = append(orders, &order)
	}

	// Any errors encountered by rows.Next or rows.Scan are returned here
	if err != nil {
		return nil, fmt.Errorf("Unable to select records: %w", err)
	}

	return orders, nil
}

func (db *Database) GetOrderByID(ctx context.Context, userId uint, id string) (*Order, error) {
	sql := `SELECT * FROM orders WHERE user_id = $1 AND id = $2`

	row := db.dbPool.QueryRow(ctx, sql, userId, id)

	var order Order
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.BondID,
		&order.Quantity,
		&order.Filled,
		&order.Price,
		&order.Action,
		&order.Status,
		&order.ExpiresAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to select record: %w", err)
	}

	return &order, nil
}

func (db *Database) GetOrdersByUserID(ctx context.Context, id uint, count uint) ([]*Order, error) {
	sql := `SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`

	rows, err := db.dbPool.Query(ctx, sql, id, count)
	if err != nil {
		return nil, fmt.Errorf("Unable to select records: %w", err)
	}
	defer rows.Close()

	orders := make([]*Order, 0)
	for rows.Next() {
		var order Order
		err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.BondID,
			&order.Quantity,
			&order.Filled,
			&order.Price,
			&order.Action,
			&order.Status,
			&order.ExpiresAt,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			break
		}
		orders = append(orders, &order)
	}

	// Any errors encountered by rows.Next or rows.Scan are returned here
	if err != nil {
		return nil, fmt.Errorf("Unable to select records: %w", err)
	}

	return orders, nil
}
