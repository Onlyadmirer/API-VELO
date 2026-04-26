package repository

import (
	"VELO-backend/internal/entity"
	"database/sql"
	"fmt"
	"log"
)

type ProductRepository interface {
	GetAllProducts() ([]entity.Product, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) GetAllProducts() ([]entity.Product, error) {
	var produccts []entity.Product

	query := "SELECT id, name, price, category, stock, image FROM products"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data products: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var p entity.Product

		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category, &p.Stock, &p.Image)
		if err != nil {
			log.Println("Error saat scan baris produk: ", err)
			continue
		}

		produccts = append(produccts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("terjadi kesalahan saat membaca baris data: %v", err)
	}

	return produccts, nil
}
