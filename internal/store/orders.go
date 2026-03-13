package store

import (
	"context"

	"github.com/your-username/shop01/internal/models"
)

func (db *DB) CreateOrder(ctx context.Context, total float64, items []models.CartItem) (int, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var orderID int
	err = tx.QueryRow(ctx,
		`INSERT INTO orders (status, total) VALUES ('paid', $1) RETURNING id`,
		total,
	).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items (order_id, product_id, name, quantity, price)
			 VALUES ($1, $2, $3, $4, $5)`,
			orderID, item.Product.ID, item.Product.Name, item.Quantity, item.Product.Price,
		)
		if err != nil {
			return 0, err
		}
	}

	return orderID, tx.Commit(ctx)
}
