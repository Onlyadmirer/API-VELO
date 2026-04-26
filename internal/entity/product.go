package entity

import "time"

type Product struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Stock     int       `json:"stock"`
	Category  string    `json:"category"`
	Price     int       `json:"price"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
