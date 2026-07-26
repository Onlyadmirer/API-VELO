package handler

import (
	"net/http"
	"os"

	"VELO-backend/pkg/service"
)

type OAuthHandler struct {
	service service.OAuthService
}

func NewOAuthHandler(service service.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		service: service,
	}
}

func (h *OAuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url, _, err := h.service.GetGoogleLoginURL()
	if err != nil {
		http.Error(w, "Gagal membuat login URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		http.Redirect(w, r, frontendURL+"/login?error=no_code", http.StatusTemporaryRedirect)
		return
	}

	cookie, _, err := h.service.HandleGoogleCallback(code, state)
	if err != nil {
		http.Redirect(w, r, frontendURL+"/login?error=oauth_failed", http.StatusTemporaryRedirect)
		return
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}
