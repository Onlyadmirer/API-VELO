package repository

import (
	"VELO-backend/internal/entity"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	CreateUser(user entity.RegisterUser) (*entity.User, error)
	FindByEmail(email string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// check existing user
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

// POST (create user)
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
