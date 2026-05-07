package handler

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/middleware"
	"VELO-backend/pkg/service"
	"encoding/json"
	"net/http"
)

// CartHandler bertanggung jawab melayani rute HTTP keranjang.
type CartHandler struct {
	service service.CartService
}

// NewCartHandler menginisialisasi instance baru untuk CartHandler.
func NewCartHandler(service service.CartService) *CartHandler {
	return &CartHandler{
		service: service,
	}
}

// POST (add product item to cart)
// AddToCart memproses request menambahkan barang ke keranjang.
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

// GET (get cart item)
// GetCart menangani request untuk melihat seluruh item pada keranjang belanja pengguna.
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(middleware.UserIdKey).(int)

	cartItem, err := h.service.GetCart(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"message": "berhasil ambil cart items", "data": cartItem})
}
