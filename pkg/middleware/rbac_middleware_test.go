package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRBACMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middlewareToTest := RBACMiddleware(nextHandler)

	// Admin
	t.Run("Admin boleh masuk", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/products/1", nil)

		ctx := context.WithValue(req.Context(), RoleKey, "admin")
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		middlewareToTest.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "gagal: Admin harusnya dapat 200 ok")
	})

	// customer
	t.Run("Customer dilarang masuk", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/products/1", nil)

		ctx := context.WithValue(req.Context(), RoleKey, "customer")
		req = req.WithContext(ctx)

		rec := httptest.NewRecorder()

		middlewareToTest.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusForbidden, rec.Code, "gagal: customer seharusnya 401")

	})

}
