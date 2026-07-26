package entity

import "time"

// User merupakan model basis data utama untuk data pengguna.
type User struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	Role          string    `json:"role"`
	IsVerified    bool      `json:"is_verified"`
	OAuthID       *string   `json:"oauth_id,omitempty"`
	OAuthProvider *string   `json:"oauth_provider,omitempty"`
	AvatarURL     *string   `json:"avatar_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type RegisterUser struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Token     string
	ExpiresAt time.Time
}
type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Token string `json:"token"`
}
