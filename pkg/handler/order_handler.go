package handler

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/middleware"
	"VELO-backend/pkg/service"
	"VELO-backend/pkg/utils"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

type OrderHandler struct {
	service service.OrderService
}

// NewOrderHandler menginisialisasi instance baru untuk OrderHandler.
func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

// CheckOut menangani proses pemesanan dengan memanggil operasi Service untuk membuat transaksi.
func (h *OrderHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idempotencyKey := r.Header.Get("X-Idempotency-Key")
	if idempotencyKey == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Idempotency key required")
		return
	}

	userIDRaw := r.Context().Value(middleware.UserIdKey)
	userID, ok := userIDRaw.(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orderID, RedirectURL, err := h.service.CreateOrder(userID, r.Context(), idempotencyKey)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Berhasil checkout", map[string]any{"order_id": orderID, "redirect_url": RedirectURL})
}

// MidtransNotifications menangani webhook/notifikasi dari Midtrans (Gateway Pembayaran).
// Berfungsi untuk membaca status pembayaran dari Payload Midtrans dan mengubah status pesanan di database.
func (h *OrderHandler) MidtransNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var notifications entity.MidtransNotifications
	if err := json.NewDecoder(r.Body).Decode(&notifications); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	serverKey := os.Getenv("SERVER_KEY")

	// validasi signature key midtrans
	rawSignature := notifications.OrderID + notifications.StatusCode + notifications.GrossAmount + serverKey

	hasher := sha512.New()
	hasher.Write([]byte(rawSignature))
	calculatedSignature := hex.EncodeToString(hasher.Sum(nil))

	if calculatedSignature != notifications.SignatureKey {

		utils.ResponseError(w, http.StatusUnauthorized, "invalid signature key")
		return
	}

	orderID, err := strconv.Atoi(notifications.OrderID)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "invalid order id")
	}

	var paymentStatus string
	switch notifications.TransactionStatus {
	case "capture", "settlement":
		paymentStatus = "Paid"
	case "cancel", "expire":
		paymentStatus = "cancel"
	default:
		paymentStatus = "Pending"
	}

	err = h.service.UpdateOrderStatus(orderID, paymentStatus)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// GetOrder menangani rute request untuk menampilkan profil dan riwayat order pelanggan.
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDRaw := r.Context().Value(middleware.UserIdKey)
	userID, ok := userIDRaw.(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	order, err := h.service.GetOrder(userID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Success", order)

}
