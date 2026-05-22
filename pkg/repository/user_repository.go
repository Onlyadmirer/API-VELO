package repository

import (
	"VELO-backend/pkg/entity"
	"database/sql"
)

// UserRepository menangani kueri ke database khusus untuk pengguna.
type UserRepository interface {
	CreateUser(user entity.RegisterUser) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	ActivateVerifyToken(token string) error
	GetUserByID(userId int) (*entity.User, error)
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

// (Register)
// CreateUser mendaftarkan user baru yang lolos verifikasi pendaftaran dan menyimpannya di DB.
func (r *userRepository) CreateUser(user entity.RegisterUser) (*entity.User, error) {
	query := `INSERT INTO users (name, email, password, verify_token, verify_token_expires_at)
	VALUES ($1, $2, $3, $4, $5) RETURNING name, email, role, created_at, updated_at`

	var u entity.User

	err := r.db.QueryRow(query, user.Name, user.Email, user.Password, user.Token, user.ExpiresAt).Scan(&u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) ActivateVerifyToken(token string) error {
	query := `
	UPDATE users SET 
	is_verified = true,
	verify_token = NULL,
	verify_token_expires_at = NULL
	WHERE verify_token = $1
	AND verify_token_expires_at > NOW()
	AND is_verified = false`

	result, err := r.db.Exec(query, token)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// (Login) accept email and return datas user for jwt auth
// GetUserByEmail mencari data pengguna beserta password yang di-hash dari database saat hendak melakukan login.
func (r *userRepository) GetUserByEmail(email string) (*entity.User, error) {
	query := `SELECT id, password, role, is_verified FROM users WHERE email = $1`
	var u entity.User
	err := r.db.QueryRow(query, email).Scan(&u.ID, &u.Password, &u.Role, &u.IsVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &u, nil
}

func (r *userRepository) GetUserByID(userId int) (*entity.User, error) {
	query := `SELECT id, name, email, role, is_verified, created_at, updated_at FROM users WHERE id = $1`

	var u entity.User
	err := r.db.QueryRow(query, userId).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Role,
		&u.IsVerified,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &u, nil
}
