package service

import (
	"VELO-backend/pkg/repository"
	"fmt"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type OrderService interface {
	CreateOrder(userId int) (int, string, error)
	UpdateOrderStatus(orderID int, status string) error
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

// create order
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
	}

	snapResp, errMidtrans := snap.CreateTransaction(resp)
	if errMidtrans != nil {
		return 0, "", fmt.Errorf("gagal membuat linkk pembayaran: %v", errMidtrans.GetMessage())
	}

	return orderID, snapResp.RedirectURL, nil

}

func (s *orderService) UpdateOrderStatus(orderID int, status string) error {
	err := s.orderRepo.UpdateOrderStatus(orderID, status)
	if err != nil {
		return err
	}
	return nil
}
