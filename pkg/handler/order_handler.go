package handler

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/middleware"
	"VELO-backend/pkg/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value(middleware.UserIdKey).(int)

	orderID, RedirectURL, err := h.service.CreateOrder(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"message": "berhasil chekout", "order_id": orderID, "redirect_url": RedirectURL})
}

// midtrans notifications
func (h *OrderHandler) MidtransNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var notifications entity.MidtransNotifications
	if err := json.NewDecoder(r.Body).Decode(&notifications); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	orderID, _ := strconv.Atoi(notifications.OrderID)

	var paymentStatus string
	switch notifications.TransactionStatus {
	case "capture", "settlement":
		paymentStatus = "Paid"
	case "cancel", "expire":
		paymentStatus = "cancel"
	default:
		paymentStatus = "Pending"
	}

	err := h.service.UpdateOrderStatus(orderID, paymentStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
