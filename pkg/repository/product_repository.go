package repository

import (
	"VELO-backend/pkg/entity"
	"database/sql"
	"fmt"
	"log"
)

type ProductRepository interface {
	GetAllProducts(page int, limit int) (entity.PaginatedProductResponse, error)
	CreateProduct(req entity.Product) error
	DeleteProduct(id int) error
	UpdateProduct(id int, req entity.Product) (*entity.Product, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

// GET
func (r *productRepository) GetAllProducts(page int, limit int) (entity.PaginatedProductResponse, error) {
	var products []entity.Product
	var totalItems int

	// total product di database
	countQuery := `SELECT COUNT(id) FROM products`
	err := r.db.QueryRow(countQuery).Scan(&totalItems)
	if err != nil {
		return entity.PaginatedProductResponse{}, fmt.Errorf("gagal hitung jumlah products")
	}

	offset := (page - 1) * limit

	query := "SELECT id, name, price, category, stock, image FROM products LIMIT $1 OFFSET $2"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return entity.PaginatedProductResponse{}, fmt.Errorf("gagal mengambil data products: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var p entity.Product

		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Stock, &p.Image)
		if err != nil {
			log.Println("Error saat scan baris produk: ", err)
			continue
		}

		products = append(products, p)
	}

	totalPages := totalItems / limit
	if totalItems%limit > 0 {

		totalPages++
	}

	if products == nil {
		products = []entity.Product{}
	}

	if err = rows.Err(); err != nil {
		return entity.PaginatedProductResponse{}, fmt.Errorf("terjadi kesalahan saat membaca baris data: %v", err)
	}

	result := entity.PaginatedProductResponse{
		Data: products,
		Metadata: entity.PaginateMeta{
			CurrentPage: page,
			TotalPages:  totalPages,
			TotalItems:  totalItems,
			Limit:       limit,
		},
	}

	return result, nil
}

// POST
func (r *productRepository) CreateProduct(req entity.Product) error {
	query := "INSERT INTO products (name, stock, category, price, image) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.db.Exec(query, req.Name, req.Stock, req.Category, req.Price, req.Image)
	if err != nil {
		return fmt.Errorf("gagal menambahkan produk: %v", err)
	}

	return nil
}

// DELETE
func (r *productRepository) DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("gagal hapus product: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("produk dengan ID %d tidak ditemukan", id)
	}

	return nil
}

// PUT
func (r *productRepository) UpdateProduct(id int, req entity.Product) (*entity.Product, error) {
	query := `UPDATE products
	SET name = $1, stock = $2, category = $3, price = $4, image = $5
	WHERE id = $6
	RETURNING id, name, stock, category, price, image`

	var product entity.Product
	err := r.db.QueryRow(query, req.Name, req.Stock, req.Category, req.Price, req.Image, id).Scan(&product.ID, &product.Name, &product.Stock, &product.Category, &product.Price, &product.Image)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("produk dengan ID %d tidak ada", id)
		}
		return nil, fmt.Errorf("gagal update data: %v", err)
	}

	return &product, nil
}
