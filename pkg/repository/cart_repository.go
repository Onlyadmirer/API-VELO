package repository

import (
	"VELO-backend/pkg/entity"
	"database/sql"
	"fmt"
)

// CartRepository mengelola operasi entri data keranjang belanja.
type CartRepository interface {
	GetOrCreateCart(userId int) (int, error)
	UpsertCartItem(cartID int, productID int, quantity int) error
	GetCart(userId int) ([]entity.CartItemResponse, error)
	UpdateCartItemQuantity(cartId int, cartItemId int, quantity int) error
	DeleteCartItem(CartId int, cartItemId int) error
	GetCartId(userID int) (int, error)
	ClearCart(userID int) error
}

type cartRepository struct {
	db *sql.DB
}

// NewCartRepository menginisialisasi implementasi CartRepository menggunakan koneksi DB.
func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{
		db: db,
	}
}

// GetOrCreateCart akan membuat entri keranjang baru jika user belum memiliki keranjang yang aktif.
func (r *cartRepository) GetOrCreateCart(userId int) (int, error) {
	query := `
		INSERT INTO carts (user_id)
		VALUES ($1)

		ON CONFLICT (user_id)
		DO UPDATE SET user_id = EXCLUDED.user_id

		RETURNING id
	`
	var id int
	err := r.db.QueryRow(query, userId).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, err
	}
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *cartRepository) GetCartId(userID int) (int, error) {
	query := `SELECT id FROM carts WHERE user_id = $1`

	var id int
	err := r.db.QueryRow(query, userID).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, err
	}

	return id, nil

}

// UpsertCartItem memperbarui jumlah (quantity) item jika sudah ada, atau menambahkannya bila belum ada.
func (r *cartRepository) UpsertCartItem(cartID int, productID int, quantity int) error {
	query := `
		INSERT INTO cart_items (
			cart_id,
			product_id,
			quantity
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET
			quantity = cart_items.quantity + EXCLUDED.quantity
	`
	_, err := r.db.Exec(query, cartID, productID, quantity)
	if err != nil {
		return fmt.Errorf("gagal tambah item ke cart: %w", err)
	}

	return nil

}

// GetCart mengambil semua baris data barang di dalam keranjang belanja milik user.
func (r *cartRepository) GetCart(userId int) ([]entity.CartItemResponse, error) {

	var cartItems []entity.CartItemResponse

	query := `
		SELECT
			ci.id,
			ci.cart_id,
			ci.quantity,
			p.id,
			p.name,
			p.price,
			p.image,
			p.category
		FROM cart_items ci
		JOIN carts c
			ON ci.cart_id = c.id
		JOIN products p
			ON ci.product_id = p.id
		WHERE c.user_id = $1
	`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("gagal query ke database: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var ci entity.CartItem

		ci.Product = &entity.Product{}

		err := rows.Scan(
			&ci.ID,
			&ci.CartID,
			&ci.Quantity,
			&ci.Product.ID,
			&ci.Product.Name,
			&ci.Product.Price,
			&ci.Product.Image,
			&ci.Product.Category,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"gagal scan cart item: %w",
				err,
			)
		}

		totalPrice := float64(ci.Quantity) * float64(ci.Product.Price)

		resp := entity.CartItemResponse{
			ID:       ci.ID,
			CartID:   ci.CartID,
			Quantity: ci.Quantity,
			Product: entity.ProductResponse{
				ID:       ci.Product.ID,
				Name:     ci.Product.Name,
				Price:    ci.Product.Price,
				Image:    ci.Product.Image,
				Category: ci.Product.Category,
			},
			TotalAmount: totalPrice,
		}

		cartItems = append(cartItems, resp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("terjadi kesalahan saat membaca baris data: %w", err)
	}

	return cartItems, nil
}

func (r *cartRepository) UpdateCartItemQuantity(cartId int, cartItemId int, quantity int) error {
	query := `UPDATE cart_items SET quantity = $1
	WHERE id = $2 AND cart_id = $3`

	result, err := r.db.Exec(query, quantity, cartItemId, cartId)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if rows == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil

}

func (r *cartRepository) DeleteCartItem(CartId int, cartItemId int) error {
	query := `DELETE FROM cart_items 
	WHERE id = $1 AND cart_id = $2`

	result, err := r.db.Exec(query, cartItemId, CartId)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

func (r *cartRepository) ClearCart(userID int) error {
	query := `DELETE FROM carts WHERE user_id = $1`

	result, err := r.db.Exec(query, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return nil
	}

	return nil
}
