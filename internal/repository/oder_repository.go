package repository

import (
	"VELO-backend/internal/entity"
	"database/sql"
	"fmt"
)

type OrderRepository interface {
	CreateOrder(userId int, cartId int, cartItems []entity.CartItemResponse) (int, float64, error)
	UpdateOrderStatus(orderID int, status string) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// membuat order/checkout
func (r *orderRepository) CreateOrder(userId int, cartId int, cartItems []entity.CartItemResponse) (orderId int, totalAmount float64, err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, 0, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// hitung total price dari cart items
	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += float64(item.Product.Price) * float64(item.Quantity)
	}

	// insert ke tabel orders
	query := `INSERT INTO orders (user_id, total_amount, status) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(query, userId, totalPrice, "Unpaid").Scan(&orderId)
	if err != nil {
		return 0, 0, fmt.Errorf("gagal insert order: %v", err)
	}

	// insert cart items ke order items
	for _, item := range cartItems {
		query := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`
		_, err = tx.Exec(query, orderId, item.Product.ID, item.Quantity, item.Product.Price)
		if err != nil {
			return 0, 0, fmt.Errorf("gagal insert order item: %v", err)
		}
	}

	// delete item yang ada di dalam cart karena sudah di order
	queryDelete := `DELETE FROM cart_items WHERE cart_id = $1`
	_, err = tx.Exec(queryDelete, cartId)
	if err != nil {
		return 0, 0, fmt.Errorf("gagal hapus item di keranjang: %v", err)
	}

	return orderId, totalPrice, nil

}

// update order status
func (r *orderRepository) UpdateOrderStatus(orderID int, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, orderID)
	if err != nil {
		return fmt.Errorf("gagal update status: %v", err)
	}

	return nil
}
