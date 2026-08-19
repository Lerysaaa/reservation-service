package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lerysaaa/reservation-service/internal/reservation"
)

type fakeRepo struct{}

func (fakeRepo) CreateProduct(context.Context, string, int) (reservation.Product, error) {
	return reservation.Product{}, nil
}

func (fakeRepo) ListProducts(context.Context) ([]reservation.Product, error) {
	return []reservation.Product{}, nil
}

func (fakeRepo) CreateReservation(context.Context, int64, int) (reservation.Reservation, error) {
	return reservation.Reservation{}, nil
}

func (fakeRepo) GetReservation(context.Context, int64) (reservation.Reservation, error) {
	return reservation.Reservation{}, reservation.ErrNotFound
}

func (fakeRepo) CancelReservation(context.Context, int64) (reservation.Reservation, error) {
	return reservation.Reservation{}, nil
}

func TestHealth(t *testing.T) {
	h := New(reservation.NewService(fakeRepo{}))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.Code)
	}
}

func TestReservationNotFound(t *testing.T) {
	h := New(reservation.NewService(fakeRepo{}))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations/10", nil)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.Code)
	}
}
