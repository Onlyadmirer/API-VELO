package handler

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/service"
	"VELO-backend/pkg/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

// ProductHandler memfasilitasi endpoint HTTP untuk produk.
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler menginisialisasi instance baru untuk ProductHandler.
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// GET
// GetAllProducts menangani endpoint yang meminta data semua produk lengkap dengan informasi paginasinya.
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	limit := 10
	if pageStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	products, err := h.service.GetAllProducts(page, limit)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Berhasil mengambil data produk", map[string]any{
		"data":     products.Data,
		"metadata": products.Metadata,
	})

}

// POST
// CreateProduct memfasilitasi pembuatan produk baru di database dari input admin.
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newProduct entity.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	err := h.service.CreateProduct(newProduct)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.ResponseSuccess(w, http.StatusOK, "Produk berhasil ditambahkan", nil)
}

// DELETE
// DeleteProduct menghapus produk secara permanen dari sistem berdasarkan param ID.
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	parseId, err := strconv.Atoi(id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	err = h.service.DeleteProduct(parseId)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Produk berhasil dihapus", nil)
}

// PUT
// UpdateProduct memproses pemutakhiran data detail suatu produk menggunakan data JSON terbaru.
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	parseId, err := strconv.Atoi(id)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var product entity.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	data, err := h.service.UpdateProduct(parseId, product)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "berhasil update product", data)

}
