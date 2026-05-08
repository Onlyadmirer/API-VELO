package handler

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/service"
	"encoding/json"
	"net/http"
)

// UserHandler menangani permintaan HTTP terkait autentikasi user.
type UserHandler struct {
	service service.UserService
}

// NewUserHandler menginisialisasi instance baru untuk UserHandler.
func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// POST (register)
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var regis entity.RegisterUser

	if err := json.NewDecoder(r.Body).Decode(&regis); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid input"})
		return
	}

	user, err := h.service.CreateUser(regis)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"message": "berhasil registrasi", "data": user})

}

// POST (Login)
func (h *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var login entity.LoginUser

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	user, err := h.service.UserLogin(login)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	http.SetCookie(w, user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"message": "Login Berhasil"})

}

func (h *UserHandler) LogOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie := &http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`"message: berhasil log out"`))
}
