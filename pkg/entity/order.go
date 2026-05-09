package entity

import "time"

type Order struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`

	CreatedAt time.Time `json:"created_at"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type MidtransNotifications struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	FraudStatus       string `json:"fraud_status"`
	GrossAmount       string `json:"gross_amount"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
}

// OrderHistory membungkus balasan riwayat pesanan.
type OrderHistory struct {
	Order OrderHistoryResponse `json:"order"`
}

type OrderHistoryResponse struct {
	ID int `json:"id"`

	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`

	CreatedAt time.Time         `json:"created_at"`
	OrderItem OrderItemResponse `json:"order_item"`
}

type OrderItemResponse struct {
	Quantity int                    `json:"quantity"`
	Product  ProductHistoryResponse `json:"product"`
}

type ProductHistoryResponse struct {
	Name string `json:"name"`
}
