package handler

import (
	"VELO-backend/internal/entity"
	"VELO-backend/internal/middleware"
	"VELO-backend/internal/service"
	"encoding/json"
	"net/http"
)

type CartHandler struct {
	service service.CartService
}

func NewCartHandler(service service.CartService) *CartHandler {
	return &CartHandler{
		service: service,
	}
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := r.Context().Value(middleware.UserIdKey).(int)

	var reqCart entity.AddCartRequest
	if err := json.NewDecoder(r.Body).Decode(&reqCart); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid input"})
		return
	}

	err := h.service.AddToCart(userID, reqCart)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "product berhasil ditambahkan ke keranjang"})
}
