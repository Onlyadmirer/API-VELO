package handler

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/middleware"
	"VELO-backend/pkg/service"
	"VELO-backend/pkg/utils"
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
		utils.ResponseError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	user, err := h.service.CreateUser(regis)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Silahkan cek email untuk verifikasi", user)

}

// GET (Handle verify token)
func (h *UserHandler) HandleVerify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.URL.Query().Get("token")
	if token == "" {
		utils.ResponseError(w, http.StatusBadRequest, "token tidak ditemukan")
		return
	}

	err := h.service.VerifyEmail(token)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Berhasil verifikasi email, silahkan kembali ke halaman login", nil)
}

// POST (Login)
func (h *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var login entity.LoginUser

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.UserLogin(login)
	if err != nil {
		utils.ResponseError(w, http.StatusUnauthorized, err.Error())
		return
	}

	http.SetCookie(w, user)

	utils.ResponseSuccess(w, http.StatusOK, "Login berhasil", nil)

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

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userIdRow := r.Context().Value(middleware.UserIdKey)

	userId, ok := userIdRow.(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dataUser, err := h.service.GetUser(userId)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
	}

	utils.ResponseSuccess(w, http.StatusOK, "Get user datas successfully", dataUser)
}
