package repository

import (
	"database/sql"
	"fmt"
)

type CartRepository interface {
	GetOrCreateCart(userId int) (int, error)
	UpsertCartItem(cartID int, productID int, quantity int) error
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
