package service

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/repository"
	"fmt"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// OrderService menangani logika transaksi order dan checkout.
type OrderService interface {
	CreateOrder(userId int) (int, string, error)
	UpdateOrderStatus(orderID int, status string) error
	GetOrder(userId int) ([]entity.OrderHistory, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
	cartRepo  repository.CartRepository
}

// NewOrderService membuat instance OrderService.
func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
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

	orderIDStr := strconv.Itoa(orderID)

	resp := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderIDStr,
			GrossAmt: int64(totalPrice),
		},
		Expiry: &snap.ExpiryDetails{
			Duration: 15,
			Unit:     "minute",
		},
	}

	snapResp, errMidtrans := snap.CreateTransaction(resp)
	if errMidtrans != nil {
		return 0, "", fmt.Errorf("gagal membuat linkk pembayaran: %v", errMidtrans.GetMessage())
	}

	return orderID, snapResp.RedirectURL, nil

}

// UpdateOrderStatus mengubah status pesanan (misal dari tertunda jadi selesai).
func (s *orderService) UpdateOrderStatus(orderID int, status string) error {
	err := s.orderRepo.UpdateOrderStatus(orderID, status)
	if err != nil {
		return err
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
