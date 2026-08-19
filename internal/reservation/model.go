package reservation

import "time"

type Product struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}

type Reservation struct {
	ID         int64      `json:"id"`
	ProductID  int64      `json:"product_id"`
	Quantity   int        `json:"quantity"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	CanceledAt *time.Time `json:"canceled_at,omitempty"`
}
