package reservation

import (
	"context"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateProduct(ctx context.Context, name string, stock int) (Product, error) {
	name = strings.TrimSpace(name)
	if name == "" || stock < 0 {
		return Product{}, ErrInvalidInput
	}
	return s.repo.CreateProduct(ctx, name, stock)
}

func (s *Service) ListProducts(ctx context.Context) ([]Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *Service) Reserve(ctx context.Context, productID int64, quantity int) (Reservation, error) {
	if productID <= 0 || quantity <= 0 {
		return Reservation{}, ErrInvalidInput
	}
	return s.repo.CreateReservation(ctx, productID, quantity)
}

func (s *Service) GetReservation(ctx context.Context, id int64) (Reservation, error) {
	if id <= 0 {
		return Reservation{}, ErrNotFound
	}
	return s.repo.GetReservation(ctx, id)
}

func (s *Service) Cancel(ctx context.Context, id int64) (Reservation, error) {
	if id <= 0 {
		return Reservation{}, ErrNotFound
	}
	return s.repo.CancelReservation(ctx, id)
}
