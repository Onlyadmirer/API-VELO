package service

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/repository"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService mendefinisikan kontrak untuk layanan pengguna (user).
type UserService interface {
	CreateUser(user entity.RegisterUser) (*entity.User, error)
	UserLogin(reqLogin entity.LoginUser) (*http.Cookie, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService membuat instance UserService baru dengan dependensi UserRepository.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// CreateUser memproses data pendaftaran, melakukan hashing password, dan menyimpan data pengguna baru ke database.
func (s *userService) CreateUser(user entity.RegisterUser) (*entity.User, error) {
	err := s.repo.FindByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal hashing password")
	}

	newUser := entity.RegisterUser{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashedPass),
	}

	dataUser, err := s.repo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return dataUser, nil

}

// UserLogin memverifikasi kredensial pengguna dan mengembalikan JWT token jika berhasil.
func (s *userService) UserLogin(reqLogin entity.LoginUser) (*http.Cookie, error) {

	dataUser, err := s.repo.GetUserByEmail(reqLogin.Email)
	if err != nil {
		return nil, err
	}

	compare := bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(reqLogin.Password))
	if compare != nil {
		return nil, fmt.Errorf("Email atau password salah")
	}

	claims := jwt.MapClaims{
		"user_id": dataUser.ID,
		"role":    dataUser.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token: %v", err)
	}

	cookie := &http.Cookie{
		Name:     "jwt_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	}

	return cookie, nil
}
