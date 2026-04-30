package service

import (
	"VELO-backend/internal/entity"
	"VELO-backend/internal/repository"
)

type CartService interface {
	AddToCart(userID int, req entity.AddCartRequest) error
	GetCart(userId int) ([]entity.CartItemResponse, error)
}

type cartService struct {
	repo repository.CartRepository
}

func NewCartService(repo repository.CartRepository) CartService {
	return &cartService{
		repo: repo,
	}
}

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

func (s *cartService) GetCart(userId int) ([]entity.CartItemResponse, error) {
	cartItem, err := s.repo.GetCart(userId)
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}
