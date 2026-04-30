package repository

import (
	"VELO-backend/internal/entity"
	"database/sql"
	"fmt"
	"log"
)

type CartRepository interface {
	GetOrCreateCart(userId int) (int, error)
	UpsertCartItem(cartID int, productID int, quantity int) error
	GetCart(userId int) ([]entity.CartItemResponse, error)
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{
		db: db,
	}
}

func (r *cartRepository) GetOrCreateCart(userId int) (int, error) {
	query := `SELECT id FROM carts WHERE user_id = $1`
	var id int
	err := r.db.QueryRow(query, userId).Scan(&id)
	if err == sql.ErrNoRows {
		queryInsert := `INSERT INTO carts (user_id) VALUES ($1) RETURNING id`
		err := r.db.QueryRow(queryInsert, userId).Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("gagal ambil id: %v", err)
		}

		return id, nil
	}

	if err != nil {
		return 0, fmt.Errorf("gagal mencari cart: %v", err)
	}

	return id, nil
}

func (r *cartRepository) UpsertCartItem(cartID int, productID int, quantity int) error {
	query := `INSERT INTO cart_items (cart_id, product_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (cart_id, product_id) DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity`
	_, err := r.db.Exec(query, cartID, productID, quantity)
	if err != nil {
		return fmt.Errorf("gagal tambah item ke cart: %v", err)
	}

	return nil

}

func (r *cartRepository) GetCart(userId int) ([]entity.CartItemResponse, error) {

	var cartItem []entity.CartItemResponse

	query := `SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity, p.name, p.price
	FROM cart_items ci
	JOIN carts c ON ci.cart_id = c.id
	JOIN products p ON ci.product_id = p.id
	WHERE c.user_id = $1`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("gagal query ke database: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var ci entity.CartItem

		ci.Product = &entity.Product{}

		if err := rows.Scan(&ci.ID, &ci.ProductID, &ci.Quantity, &ci.Product.ID, &ci.Product.Name, &ci.Product.Price); err != nil {
			log.Println("error saat scan baris cart items: ", err)
			continue
		}

		resp := entity.CartItemResponse{
			ID:       ci.ID,
			Quantity: ci.Quantity,
			Product: entity.ProductResponse{
				ID:    ci.Product.ID,
				Name:  ci.Product.Name,
				Price: ci.Product.Price,
			},
		}

		cartItem = append(cartItem, resp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("terjadi kesalahan saat membaca baris data: %v", err)
	}

	return cartItem, nil
}
