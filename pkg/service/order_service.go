package service

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/payment"
	"VELO-backend/pkg/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// OrderService menangani logika transaksi order dan checkout.
type OrderService interface {
	CreateOrder(userId int, ctx context.Context, idempotencyKey string) (int, string, error)
	UpdateOrderStatus(orderID int, status string) error
	GetOrder(userId int) ([]entity.OrderHistory, error)
	CancelExpiredOrders()
}

type orderService struct {
	orderRepo      repository.OrderRepository
	cartRepo       repository.CartRepository
	PaymentGateway payment.PaymentGateway
	redis          *redis.Client
}

// NewOrderService membuat instance OrderService.
func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, paymentGateway payment.PaymentGateway, redis *redis.Client) OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		cartRepo:       cartRepo,
		PaymentGateway: paymentGateway,
		redis:          redis,
	}
}

// CreateOrder menghitung harga, membuat pesanan di basis data, dan memulai proses pembayaran dengan payment gateway.
func (s *orderService) CreateOrder(userId int, ctx context.Context, idempotencyKey string) (int, string, error) {

	redisKey := fmt.Sprintf("velo:idempotency:%s", idempotencyKey)

	//  cek apakah di redis ada idempotency key
	if s.redis != nil {
		result, err := s.redis.Get(ctx, redisKey).Result()
		if err == nil {
			var cached struct {
				OrderId     int    `json:"order_id"`
				RedirectUrl string `json:"redirect_url"`
			}
			if err := json.Unmarshal([]byte(result), &cached); err != nil {
				fmt.Println("cache corrupted, reprocess:", err)
			} else {

				return cached.OrderId, cached.RedirectUrl, nil
			}
		}
		if err != redis.Nil {
			fmt.Println("koneksi redis bermasalah: ", err)
		}
	}

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

	// masukkan ke redis
	if s.redis != nil {
		cacheData, _ := json.Marshal(map[string]any{
			"order_id":     orderID,
			"redirect_url": redirectURL,
		})
		s.redis.Set(ctx, redisKey, cacheData, 15*time.Minute)
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

// CancelExpiredOrders menemukan order Unpaid yang sudah expired dan membatalkannya.
func (s *orderService) CancelExpiredOrders() {
	const expiryMinutes = 15

	orderIDs, err := s.orderRepo.GetExpiredUnpaidOrders(expiryMinutes)
	if err != nil {
		fmt.Println("gagal mengambil expired orders:", err)
		return
	}

	for _, id := range orderIDs {
		if err := s.UpdateOrderStatus(id, entity.OrderStatusCancel); err != nil {
			fmt.Printf("gagal membatalkan expired order %d: %v\n", id, err)
		}
	}
}

// GetOrder meretrieve riwayat pesanan beserta status pembayarannya.
func (s *orderService) GetOrder(userId int) ([]entity.OrderHistory, error) {
	order, err := s.orderRepo.GetOrder(userId)
	if err != nil {
		return nil, err
	}

	return order, nil
}
