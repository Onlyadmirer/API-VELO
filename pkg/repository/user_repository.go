package repository

import (
	"VELO-backend/pkg/entity"
	"database/sql"
	"fmt"
)

// UserRepository menangani kueri ke database khusus untuk pengguna.
type UserRepository interface {
	CreateUser(user entity.RegisterUser) (*entity.User, error)
	FindByEmail(email string) error
	GetUserByEmail(email string) (*entity.User, error)
	// UserLogin(user entity.LoginUser) (*entity.LoginResponse, error)
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository menginisialisasi implementasi UserRepository menggunakan koneksi DB.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// check existing user
// FindByEmail memvalidasi kepemilikan alamat email (untuk mencegah duplikasi pendaftaran email).
func (r *userRepository) FindByEmail(email string) error {
	query := `SELECT id, email, password, role FROM users WHERE email = $1`
	var user entity.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	return fmt.Errorf("email sudah terdaftar")
}

// accept email and return datas user for jwt auth
// GetUserByEmail mencari data pengguna beserta password yang di-hash dari database saat hendak melakukan login.
func (r *userRepository) GetUserByEmail(email string) (*entity.User, error) {
	query := `SELECT id, password, role FROM users WHERE email = $1`
	var u entity.User
	err := r.db.QueryRow(query, email).Scan(&u.ID, &u.Password, &u.Role)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil data user")
	}

	return &u, nil
}

// (Register)
// CreateUser mendaftarkan user baru yang lolos verifikasi pendaftaran dan menyimpannya di DB.
func (r *userRepository) CreateUser(user entity.RegisterUser) (*entity.User, error) {
	query := `INSERT INTO users (name, email, password)
	VALUES ($1, $2, $3) RETURNING name, role, created_at, updated_at`

	var u entity.User

	err := r.db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("gagal create user: %v", err)
	}

	return &u, nil
}
