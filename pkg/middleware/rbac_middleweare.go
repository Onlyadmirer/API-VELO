package middleware

import (
	"encoding/json"
	"net/http"
)

func RBACMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userRole, ok := r.Context().Value(RoleKey).(string)
		if !ok || userRole != "admin" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Akses ditolak, khusus admin"})
			return
		}

		next.ServeHTTP(w, r)
	}
}
