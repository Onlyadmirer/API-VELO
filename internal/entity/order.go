package entity

type Order struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
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
}
