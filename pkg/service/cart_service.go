package service

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/repository"
)

// CartService mendefinisikan kontrak untuk layanan keranjang belanja pelanggan.
type CartService interface {
	AddToCart(userID int, req entity.AddCartRequest) error
	GetCart(userId int) ([]entity.CartItemResponse, error)
}

type cartService struct {
	repo repository.CartRepository
}

// NewCartService membuat instance CartService.
func NewCartService(repo repository.CartRepository) CartService {
	return &cartService{
		repo: repo,
	}
}

// AddToCart memasukkan barang ke keranjang atau menambah kuantitasnya.
func (s *cartService) AddToCart(userID int, req entity.AddCartRequest) error {
	id, err := s.repo.GetOrCreateCart(userID)
	if err != nil {
		return err
	}

	err = s.repo.UpsertCartItem(id, req.ProductID, req.Quantity)
	if err != nil {
		return err
	}

	return nil

}

// GetCart mengambil semua item di dalam keranjang pengguna.
func (s *cartService) GetCart(userId int) ([]entity.CartItemResponse, error) {
	cartItem, err := s.repo.GetCart(userId)
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}
