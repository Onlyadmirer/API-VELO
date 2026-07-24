package entity

import "time"

// Cart merepresentasikan data keranjang belanja milik user.
type Cart struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CartItem struct {
	ID        int       `json:"id"`
	CartID    int       `json:"cart_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Product   *Product  `json:"product,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AddCartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CartItemResponse struct {
	ID          int             `json:"id"`
	CartID      int             `json:"cart_id"`
	Quantity    int             `json:"quantity"`
	Product     ProductResponse `json:"product"`
	TotalAmount float64         `json:"total_amount"`
}

type ProductResponse struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Image    string  `json:"image"`
	Category string  `json:"category"`
	Stock    int     `json:"stock"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity"`
}
