package reservation

import (
	"context"
	"errors"
	"testing"
)

type stubRepository struct {
	createdProduct Product
	createdRes     Reservation
	createResErr   error
	cancelErr      error
}

func (s *stubRepository) CreateProduct(context.Context, string, int) (Product, error) {
	return s.createdProduct, nil
}

func (s *stubRepository) ListProducts(context.Context) ([]Product, error) {
	return nil, nil
}

func (s *stubRepository) CreateReservation(context.Context, int64, int) (Reservation, error) {
	return s.createdRes, s.createResErr
}

func (s *stubRepository) GetReservation(context.Context, int64) (Reservation, error) {
	return s.createdRes, nil
}

func (s *stubRepository) CancelReservation(context.Context, int64) (Reservation, error) {
	return s.createdRes, s.cancelErr
}

func TestReserveRejectsInvalidQuantity(t *testing.T) {
	svc := NewService(&stubRepository{})
	_, err := svc.Reserve(context.Background(), 1, 0)
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestReservePassesRepositoryError(t *testing.T) {
	repo := &stubRepository{createResErr: ErrInsufficientStock}
	svc := NewService(repo)
	_, err := svc.Reserve(context.Background(), 1, 5)
	if !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
}
