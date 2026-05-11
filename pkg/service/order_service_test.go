package service

import (
	"VELO-backend/pkg/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCartRepo struct {
	mock.Mock
}
type MockOrderRepo struct {
	mock.Mock
}

func (m *MockCartRepo) GetCart(userId int) ([]entity.CartItemResponse, error) {
	args := m.Called(userId)
	return args.Get(0).([]entity.CartItemResponse), args.Error(1)
}

func (m *MockCartRepo) GetOrCreateCart(userId int) (int, error)                      { return 0, nil }
func (m *MockCartRepo) UpsertCartItem(cartID int, productID int, quantity int) error { return nil }

func (m *MockOrderRepo) CreateOrder(u int, c int, items []entity.CartItemResponse) (int, float64, error) {
	return 0, 0, nil
}

func (m *MockOrderRepo) UpdateOrderStatus(orderID int, status string) error { return nil }
func (m *MockOrderRepo) GetOrder(userId int) ([]entity.OrderHistory, error) { return nil, nil }
func (m *MockOrderRepo) RestoreStock(orderID int) error                     { return nil }
func (m *MockOrderRepo) GetOrderStatus(orderID int) (string, error)         { return "", nil }

func TestOrderService(t *testing.T) {

	t.Run("cart kosong", func(t *testing.T) {
		userID := 99

		cartRepoFake := new(MockCartRepo)
		orderRepoFake := new(MockOrderRepo)

		orderService := NewOrderService(orderRepoFake, cartRepoFake)

		cartRepoFake.On("GetCart", userID).Return([]entity.CartItemResponse{}, nil)

		orderID, redirectURL, err := orderService.CreateOrder(userID)

		assert.Error(t, err)
		assert.Equal(t, "keranjang masih kosong", err.Error())
		assert.Equal(t, 0, orderID)
		assert.Equal(t, "", redirectURL)
	})
}
