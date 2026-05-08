package repository

import (
	"VELO-backend/pkg/entity"
	"database/sql"
	"fmt"
	"log"
)

type OrderRepository interface {
	CreateOrder(userId int, cartId int, cartItems []entity.CartItemResponse) (int, float64, error)
	UpdateOrderStatus(orderID int, status string) error
	GetOrder(userId int) ([]entity.OrderHistory, error)
}

type orderRepository struct {
	db *sql.DB
}

// NewOrderRepository menginisialisasi implementasi OrderRepository dengan koneksi SQL DB.
func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// CreateOrder mengeksekusi transaksi DB: Hitung harga -> Insert 'orders' -> Insert 'order_items' -> Kosongkan keranjang.
// Semua query dieksekusi di DALAM transaksi, di-rollback secara otomatis pakai 'defer' jika terjadi error/panic.
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

		// kurangi stock product
		stock := item.Quantity

		queryKurangiStock := `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`
		_, err = tx.Exec(queryKurangiStock, stock, item.Product.ID)
		if err != nil {
			return 0, 0, fmt.Errorf("gagla kurangi stock product: %v", err)
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

// UpdateOrderStatus mengubah status pesanan yang ada di database (contoh: Unpaid -> Paid).
func (r *orderRepository) UpdateOrderStatus(orderID int, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, orderID)
	if err != nil {
		return fmt.Errorf("gagal update status: %v", err)
	}

	return nil
}

// GetOrder meretrieve riwayat pesanan user menggunakan JOIN antara sql table orders, order_items, dan products.
func (r *orderRepository) GetOrder(userId int) ([]entity.OrderHistory, error) {

	var orderHistory []entity.OrderHistory

	query := `SELECT oi.order_id, oi.quantity, o.total_amount, o.status, o.created_at, p.name
	FROM order_items oi
	JOIN orders o ON oi.order_id = o.id
	JOIN products p ON oi.product_id = p.id
	WHERE o.user_id = $1`
	order, err := r.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("gagal query ke database: %v", err)
	}

	defer order.Close()

	for order.Next() {
		var ord entity.OrderHistory

		if err := order.Scan(&ord.Order.ID, &ord.Order.OrderItem.Quantity, &ord.Order.TotalAmount, &ord.Order.Status, &ord.Order.CreatedAt, &ord.Order.OrderItem.Product.Name); err != nil {
			log.Println("error saat scan baris order history: ", err)
			continue
		}

		resp := entity.OrderHistory{

			Order: entity.OrderHistoryResponse{
				ID:          ord.Order.ID,
				TotalAmount: ord.Order.TotalAmount,
				Status:      ord.Order.Status,
				CreatedAt:   ord.Order.CreatedAt,
				OrderItem: entity.OrderItemResponse{
					Quantity: ord.Order.OrderItem.Quantity,
					Product: entity.ProductHistoryResponse{
						Name: ord.Order.OrderItem.Product.Name,
					},
				},
			},
		}

		orderHistory = append(orderHistory, resp)
	}

	if err = order.Err(); err != nil {
		return nil, fmt.Errorf("terjadi kesalahan saat membaca baris data: %v", err)
	}

	return orderHistory, nil
}
