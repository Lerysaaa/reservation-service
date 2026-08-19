package postgres

import (
	"context"
	"errors"

	"github.com/Lerysaaa/reservation-service/internal/reservation"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateProduct(ctx context.Context, name string, stock int) (reservation.Product, error) {
	var product reservation.Product
	err := r.pool.QueryRow(ctx,
		`INSERT INTO products (name, stock) VALUES ($1, $2)
		 RETURNING id, name, stock, created_at`,
		name, stock,
	).Scan(&product.ID, &product.Name, &product.Stock, &product.CreatedAt)
	return product, err
}

func (r *Repository) ListProducts(ctx context.Context) ([]reservation.Product, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, stock, created_at FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]reservation.Product, 0)
	for rows.Next() {
		var product reservation.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Stock, &product.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (r *Repository) CreateReservation(ctx context.Context, productID int64, quantity int) (reservation.Reservation, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return reservation.Reservation{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var stock int
	err = tx.QueryRow(ctx, `SELECT stock FROM products WHERE id = $1 FOR UPDATE`, productID).Scan(&stock)
	if errors.Is(err, pgx.ErrNoRows) {
		return reservation.Reservation{}, reservation.ErrNotFound
	}
	if err != nil {
		return reservation.Reservation{}, err
	}
	if stock < quantity {
		return reservation.Reservation{}, reservation.ErrInsufficientStock
	}

	if _, err = tx.Exec(ctx, `UPDATE products SET stock = stock - $1 WHERE id = $2`, quantity, productID); err != nil {
		return reservation.Reservation{}, err
	}

	var result reservation.Reservation
	err = tx.QueryRow(ctx,
		`INSERT INTO reservations (product_id, quantity)
		 VALUES ($1, $2)
		 RETURNING id, product_id, quantity, status, created_at, canceled_at`,
		productID, quantity,
	).Scan(&result.ID, &result.ProductID, &result.Quantity, &result.Status, &result.CreatedAt, &result.CanceledAt)
	if err != nil {
		return reservation.Reservation{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return reservation.Reservation{}, err
	}
	return result, nil
}

func (r *Repository) GetReservation(ctx context.Context, id int64) (reservation.Reservation, error) {
	var result reservation.Reservation
	err := r.pool.QueryRow(ctx,
		`SELECT id, product_id, quantity, status, created_at, canceled_at
		 FROM reservations WHERE id = $1`, id,
	).Scan(&result.ID, &result.ProductID, &result.Quantity, &result.Status, &result.CreatedAt, &result.CanceledAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return reservation.Reservation{}, reservation.ErrNotFound
	}
	return result, err
}

func (r *Repository) CancelReservation(ctx context.Context, id int64) (reservation.Reservation, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return reservation.Reservation{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var result reservation.Reservation
	err = tx.QueryRow(ctx,
		`SELECT id, product_id, quantity, status, created_at, canceled_at
		 FROM reservations WHERE id = $1 FOR UPDATE`, id,
	).Scan(&result.ID, &result.ProductID, &result.Quantity, &result.Status, &result.CreatedAt, &result.CanceledAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return reservation.Reservation{}, reservation.ErrNotFound
	}
	if err != nil {
		return reservation.Reservation{}, err
	}
	if result.Status == "canceled" {
		return reservation.Reservation{}, reservation.ErrAlreadyCanceled
	}

	if _, err = tx.Exec(ctx, `UPDATE products SET stock = stock + $1 WHERE id = $2`, result.Quantity, result.ProductID); err != nil {
		return reservation.Reservation{}, err
	}

	err = tx.QueryRow(ctx,
		`UPDATE reservations SET status = 'canceled', canceled_at = now()
		 WHERE id = $1
		 RETURNING id, product_id, quantity, status, created_at, canceled_at`, id,
	).Scan(&result.ID, &result.ProductID, &result.Quantity, &result.Status, &result.CreatedAt, &result.CanceledAt)
	if err != nil {
		return reservation.Reservation{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return reservation.Reservation{}, err
	}
	return result, nil
}
