package service

import (
	"VELO-backend/internal/repository"
	"fmt"
)

type OrderService interface {
	CreateOrder(userId int) (int, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
	cartRepo  repository.CartRepository
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
	}
}

func (s *orderService) CreateOrder(userId int) (int, error) {
	cartItems, err := s.cartRepo.GetCart(userId)
	if err != nil {
		return 0, err
	}

	if len(cartItems) <= 0 {
		return 0, fmt.Errorf("keranjang masih kosong")
	}

	cartId := cartItems[0].CartID

	orderID, err := s.orderRepo.CreateOrder(userId, cartId, cartItems)
	if err != nil {
		return 0, err
	}

	return orderID, nil

}
