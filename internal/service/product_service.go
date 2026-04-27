package service

import (
	"VELO-backend/internal/entity"
	"VELO-backend/internal/repository"
	"errors"
)

type ProductService interface {
	GetAllProducts() ([]entity.Product, error)
	CreateProduct(req entity.Product) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) GetAllProducts() ([]entity.Product, error) {

	products, err := s.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}

	return products, nil
}

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
