package handler

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/middleware"
	"VELO-backend/pkg/service"
	"VELO-backend/pkg/utils"
	"encoding/json"
	"net/http"
	"strconv"
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
	userIDRaw := r.Context().Value(middleware.UserIdKey)
	userID, ok := userIDRaw.(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var reqCart entity.AddCartRequest
	if err := json.NewDecoder(r.Body).Decode(&reqCart); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid input: "+err.Error())
		return
	}

	err := h.service.AddToCart(userID, reqCart)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "product berhasil ditambahkan ke keranjang", nil)
}

// GET (get cart item)
// GetCart menangani request untuk melihat seluruh item pada keranjang belanja pengguna.
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDRaw := r.Context().Value(middleware.UserIdKey)
	userID, ok := userIDRaw.(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cartItem, err := h.service.GetCart(userID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "berhasil ambil cart items", cartItem)
}

func (h *CartHandler) UpdateCartItemQuantity(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIdRaw := r.Context().Value(middleware.UserIdKey)
	userId, ok := userIdRaw.(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")

	var quantityUpdate entity.UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&quantityUpdate); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "input invalid")
		return
	}

	cartItemIdParse, err := strconv.Atoi(id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid ID")
		return
	}

	cartItem, err := h.service.UpdateCartItemQuantity(userId, cartItemIdParse, quantityUpdate.Quantity)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "update successfully", cartItem)
}
