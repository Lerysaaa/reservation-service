package reservation

import "context"

type Repository interface {
	CreateProduct(ctx context.Context, name string, stock int) (Product, error)
	ListProducts(ctx context.Context) ([]Product, error)
	CreateReservation(ctx context.Context, productID int64, quantity int) (Reservation, error)
	GetReservation(ctx context.Context, id int64) (Reservation, error)
	CancelReservation(ctx context.Context, id int64) (Reservation, error)
}
