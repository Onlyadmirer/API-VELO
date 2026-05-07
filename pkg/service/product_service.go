package service

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/repository"
	"errors"
)

// ProductService mendefinisikan kontrak untuk layanan produk.
type ProductService interface {
	GetAllProducts(page int, limit int) (entity.PaginatedProductResponse, error)
	CreateProduct(req entity.Product) error
	DeleteProduct(id int) error
	UpdateProduct(id int, req entity.Product) (*entity.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

// NewProductService membuat instance ProductService baru.
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

// GetAllProducts mengambil daftar produk dengan mendukung sistem paginasi.
func (s *productService) GetAllProducts(page int, limit int) (entity.PaginatedProductResponse, error) {

	products, err := s.repo.GetAllProducts(page, limit)
	if err != nil {
		return entity.PaginatedProductResponse{}, err
	}

	return products, nil
}

// CreateProduct menangani pembuatan data produk baru dan validasi awal.
func (s *productService) CreateProduct(req entity.Product) error {

	if req.Name == "" {
		return errors.New("nama produk tidak boleh kosong")
	}

	if req.Price <= 0 {
		return errors.New("price harus lebih dari 0")
	}

	if req.Stock <= 0 {
		return errors.New("stock harus lebih dari 0")
	}

	err := s.repo.CreateProduct(req)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProduct menghapus produk berdasarkan ID dari database.
func (s *productService) DeleteProduct(id int) error {
	err := s.repo.DeleteProduct(id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProduct memperbarui data produk berdasarkan ID.
func (s *productService) UpdateProduct(id int, req entity.Product) (*entity.Product, error) {
	if req.Name == "" {
		return nil, errors.New("nama produk tidak boleh kosong")
	}

	if req.Price <= 0 {
		return nil, errors.New("price harus lebih dari 0")
	}

	if req.Stock <= 0 {
		return nil, errors.New("stock harus lebih dari 0")
	}

	data, err := s.repo.UpdateProduct(id, req)
	if err != nil {
		return nil, err
	}

	return data, nil
}
