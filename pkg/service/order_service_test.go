package service

import (
	"VELO-backend/pkg/entity"
	"errors"
	"os"
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
type MockPaymentGateway struct {
	mock.Mock
}

func (m *MockPaymentGateway) GenerateSnapURL(orderID int, totalPrice float64) (string, error) {
	args := m.Called(orderID, totalPrice)
	return args.String(0), args.Error(1)
}

func (m *MockCartRepo) GetCart(userId int) ([]entity.CartItemResponse, error) {
	args := m.Called(userId)
	return args.Get(0).([]entity.CartItemResponse), args.Error(1)
}

func (m *MockOrderRepo) CreateOrder(u int, c int, items []entity.CartItemResponse) (int, float64, error) {
	args := m.Called(u, c, items)
	return args.Int(0), args.Get(1).(float64), args.Error(2)
}

func (m *MockCartRepo) GetOrCreateCart(userId int) (int, error)                      { return 0, nil }
func (m *MockCartRepo) UpsertCartItem(cartID int, productID int, quantity int) error { return nil }

func (m *MockOrderRepo) UpdateOrderStatus(orderID int, status string) error { return nil }
func (m *MockOrderRepo) GetOrder(userId int) ([]entity.OrderHistory, error) { return nil, nil }
func (m *MockOrderRepo) RestoreStock(orderID int) error                     { return nil }
func (m *MockOrderRepo) GetOrderStatus(orderID int) (string, error)         { return "", nil }

func TestOrderService(t *testing.T) {

	os.Setenv("SERVER_KEY", "midtrans-test")
	defer os.Unsetenv("SERVER_KEY")

	t.Run("cart kosong", func(t *testing.T) {
		userID := 99

		cartRepoFake := new(MockCartRepo)
		orderRepoFake := new(MockOrderRepo)
		midtransFake := new(MockPaymentGateway)

		orderService := NewOrderService(orderRepoFake, cartRepoFake, midtransFake)

		cartRepoFake.On("GetCart", userID).Return([]entity.CartItemResponse{}, nil)

		orderID, redirectURL, err := orderService.CreateOrder(userID)

		assert.Error(t, err)
		assert.Equal(t, "keranjang masih kosong", err.Error())
		assert.Equal(t, 0, orderID)
		assert.Equal(t, "", redirectURL)
	})

	t.Run("database error", func(t *testing.T) {
		userID := 99

		cartRepoFake := new(MockCartRepo)
		orderRepoFake := new(MockOrderRepo)
		midtransFake := new(MockPaymentGateway)

		orderService := NewOrderService(orderRepoFake, cartRepoFake, midtransFake)

		dbErr := errors.New("database timeout")

		cartRepoFake.On("GetCart", userID).Return([]entity.CartItemResponse{}, dbErr)

		orderID, redirectURL, err := orderService.CreateOrder(userID)

		assert.Error(t, err)
		assert.Equal(t, "database timeout", err.Error())
		assert.Equal(t, 0, orderID)
		assert.Equal(t, "", redirectURL)

		cartRepoFake.AssertExpectations(t)
	})

	t.Run("checkout success", func(t *testing.T) {

		userID := 99

		cartRepoFake := new(MockCartRepo)
		orderRepoFake := new(MockOrderRepo)
		midtransFake := new(MockPaymentGateway)

		orderService := NewOrderService(orderRepoFake, cartRepoFake, midtransFake)

		mockItems := []entity.CartItemResponse{
			{ID: 1, Quantity: 2},
		}

		cartRepoFake.On("GetCart", userID).Return(mockItems, nil)

		orderRepoFake.On("CreateOrder", userID, mock.Anything, mock.Anything).Return(101, float64(500000), nil)

		midtransFake.On("GenerateSnapURL", mock.Anything, mock.Anything).Return("url-test.com", nil)

		orderID, redirectURL, err := orderService.CreateOrder(userID)

		assert.NoError(t, err)
		assert.Equal(t, 101, orderID)
		assert.NotEmpty(t, redirectURL)

		cartRepoFake.AssertExpectations(t)
		orderRepoFake.AssertExpectations(t)
		midtransFake.AssertExpectations(t)
	})
}
