package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateTestToken(userId int, role string, secret string) string {
	claims := jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		panic("gagal membuat token test" + err.Error())
	}
	return signedToken
}

func TestJWTMiddleware(t *testing.T) {
	os.Setenv("SECRET_KEY", "rahasia-negara")

	defer os.Unsetenv("SECRET_KEY")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middlewareToTest := JWTMiddleware(nextHandler)

	t.Run("tidak ada token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/cart", nil)

		rec := httptest.NewRecorder()

		middlewareToTest.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

	})

	t.Run("token salah", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/cart", nil)
		fakeToken := generateTestToken(1, "customer", "rahasia-jokowi")
		cookie := &http.Cookie{
			Name:  "jwt_token",
			Value: fakeToken,
		}

		req.AddCookie(cookie)
		rec := httptest.NewRecorder()

		middlewareToTest.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code, "unauthorized")

	})

	t.Run("token valid", func(t *testing.T) {
		token := generateTestToken(1, "customer", "rahasia-negara")
		cookie := &http.Cookie{
			Name:  "jwt_token",
			Value: token,
		}

		req := httptest.NewRequest(http.MethodPost, "/api/cart", nil)

		req.AddCookie(cookie)
		rec := httptest.NewRecorder()

		middlewareToTest.ServeHTTP(rec, req)
		fmt.Println("Pesan Error dari Middleware:", rec.Body.String())

		assert.Equal(t, http.StatusOK, rec.Code)
	})

}
