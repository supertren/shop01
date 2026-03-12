package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/your-username/shop01/internal/models"
)

var ErrNotFound = errors.New("record not found")

func (db *DB) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, image_url, stock
		FROM products
		WHERE id = $1
	`
	row := db.Pool.QueryRow(ctx, query, id)

	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.Stock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (db *DB) ListProducts(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT id, name, description, price, image_url, stock
		FROM products
		ORDER BY id ASC
	`
	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows)
}

func (db *DB) ListFeaturedProducts(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT id, name, description, price, image_url, stock
		FROM products
		WHERE stock > 0
		ORDER BY id ASC
		LIMIT 4
	`
	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProducts(rows)
}

// scanProducts iterates over rows and scans them into a slice of Products.
func scanProducts(rows pgx.Rows) ([]models.Product, error) {
	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}
