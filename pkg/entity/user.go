package entity

import "time"

// User merupakan model basis data utama untuk data pengguna.
type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	Role       string `json:"role"`
	IsVerified bool
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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
