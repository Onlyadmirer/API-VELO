package service

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/payment"
	"VELO-backend/pkg/repository"
	"fmt"
)

// OrderService menangani logika transaksi order dan checkout.
type OrderService interface {
	CreateOrder(userId int) (int, string, error)
	UpdateOrderStatus(orderID int, status string) error
	GetOrder(userId int) ([]entity.OrderHistory, error)
}

type orderService struct {
	orderRepo      repository.OrderRepository
	cartRepo       repository.CartRepository
	PaymentGateway payment.PaymentGateway
}

// NewOrderService membuat instance OrderService.
func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, paymentGateway payment.PaymentGateway) OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		cartRepo:       cartRepo,
		PaymentGateway: paymentGateway,
	}
}

// CreateOrder menghitung harga, membuat pesanan di basis data, dan memulai proses pembayaran dengan payment gateway.
func (s *orderService) CreateOrder(userId int) (int, string, error) {
	cartItems, err := s.cartRepo.GetCart(userId)
	if err != nil {
		return 0, "", err
	}

	if len(cartItems) <= 0 {
		return 0, "", fmt.Errorf("keranjang masih kosong")
	}

	cartId := cartItems[0].CartID

	orderID, totalPrice, err := s.orderRepo.CreateOrder(userId, cartId, cartItems)
	if err != nil {
		return 0, "", err
	}

	redirectURL, err := s.PaymentGateway.GenerateSnapURL(orderID, totalPrice)
	if err != nil {
		return 0, "", err
	}

	return orderID, redirectURL, nil

}

// UpdateOrderStatus mengubah status pesanan (misal dari tertunda jadi selesai).
func (s *orderService) UpdateOrderStatus(orderID int, status string) error {

	currentStatus, err := s.orderRepo.GetOrderStatus(orderID)
	if err != nil {
		return err
	}

	if currentStatus == "cancel" {
		return nil
	}

	err = s.orderRepo.UpdateOrderStatus(orderID, status)
	if err != nil {
		return err
	}

	if status == "cancel" {
		errRestock := s.orderRepo.RestoreStock(orderID)
		if errRestock != nil {
			return errRestock
		}
	}
	return nil
}

// GetOrder meretrieve riwayat pesanan beserta status pembayarannya.
func (s *orderService) GetOrder(userId int) ([]entity.OrderHistory, error) {
	order, err := s.orderRepo.GetOrder(userId)
	if err != nil {
		return nil, err
	}

	return order, nil
}
