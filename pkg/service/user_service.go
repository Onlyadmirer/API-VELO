package service

import (
	"VELO-backend/pkg/entity"
	"VELO-backend/pkg/helper"
	"VELO-backend/pkg/repository"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserService mendefinisikan kontrak untuk layanan pengguna (user).
type UserService interface {
	CreateUser(user entity.RegisterUser) (*entity.User, error)
	UserLogin(reqLogin entity.LoginUser) (*http.Cookie, error)
	VerifyEmail(token string) error
}

type userService struct {
	repo  repository.UserRepository
	email EmailService
}

// NewUserService membuat instance UserService baru dengan dependensi UserRepository.
func NewUserService(repo repository.UserRepository, email EmailService) UserService {
	return &userService{
		repo:  repo,
		email: email,
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

	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	newUser := entity.RegisterUser{
		Name:      user.Name,
		Email:     user.Email,
		Password:  string(hashedPass),
		Token:     token,
		ExpiresAt: expiresAt,
	}

	dataUser, err := s.repo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	go func() {
		errEmail := s.email.SendVerificationEmail(user.Email, user.Name, token)
		if errEmail != nil {
			log.Println("[URGENT] Gagal kirim email verifikasi:", errEmail)
		}
	}()

	return dataUser, nil
}

func (s *userService) VerifyEmail(token string) error {
	err := s.repo.FindByVerifyToken(token)
	if err != nil {
		return err
	}

	return nil
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

	if dataUser.IsVerified != true {
		return nil, fmt.Errorf("Email belum di verifikasi")
	}

	jwtToken, err := helper.GenerateJWTToken(dataUser.ID, dataUser.Role)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "jwt_token",
		Value:    jwtToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
	}

	return cookie, nil
}
