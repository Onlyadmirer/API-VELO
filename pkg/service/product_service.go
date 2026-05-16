package service

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// ProductService mendefinisikan kontrak untuk layanan produk.
type ProductService interface {
	GetAllProducts(tx context.Context, page int, limit int) (entity.PaginatedProductResponse, error)
	CreateProduct(req entity.Product) error
	DeleteProduct(id int) error
	UpdateProduct(id int, req entity.Product) (*entity.Product, error)
}

type productService struct {
	repo  repository.ProductRepository
	redis *redis.Client
}

// NewProductService membuat instance ProductService baru.
func NewProductService(repo repository.ProductRepository, redis *redis.Client) ProductService {
	return &productService{
		repo:  repo,
		redis: redis,
	}
}

// GetAllProducts mengambil daftar produk dengan mendukung sistem paginasi.
func (s *productService) GetAllProducts(ctx context.Context, page int, limit int) (entity.PaginatedProductResponse, error) {

	cacheKey := fmt.Sprintf("velo:products:page:%d:limit%d", page, limit)

	cachedData, err := s.redis.Get(ctx, cacheKey).Result()

	if err != nil && err != redis.Nil {
		fmt.Println("koneksi redis bermasalah: ", err)
	}

	if err == nil {
		log.Printf("CACHE HIT: untuk kunci: %s\n", cacheKey)

		var products entity.PaginatedProductResponse
		if json.Unmarshal([]byte(cachedData), &products) == nil {
			return products, nil
		}
	}

	// jika cache kosong maka ambil di db
	products, err := s.repo.GetAllProducts(page, limit)
	if err != nil {
		return entity.PaginatedProductResponse{}, err
	}

	// memasukkan ke cache
	jsonData, err := json.Marshal(products)
	if err == nil {
		_ = s.redis.Set(ctx, cacheKey, jsonData, 5*time.Minute).Err()
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
