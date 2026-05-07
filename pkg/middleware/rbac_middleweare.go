package middleware

import (
	"encoding/json"
	"net/http"
)

func RBACMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Context().Value(RoleKey)

		value, ok := userRole.(string)
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Akses ditolak, khusus admin"})
			return
		}

		if value == "admin" {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Akses ditolak, khusus admin"})
			return
		}
	}
}
