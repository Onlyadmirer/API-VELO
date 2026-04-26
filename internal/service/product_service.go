package service

import (
	"VELO-backend/internal/entity"
	"VELO-backend/internal/repository"
)

type ProductService interface {
	GetAllProducts() ([]entity.Product, error)
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
