package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ulbithebest/BE-pendaftaran/internal/auth"
)

// PasetoPayloadKey adalah tipe custom untuk context key
type PasetoPayloadKey string

const payloadKey PasetoPayloadKey = "pasetoPayload"

// AuthMiddleware melindungi endpoint
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Authorization header required"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error": "Invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		payload, err := auth.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, `{"error": "Invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Simpan payload di context untuk digunakan di handler selanjutnya
		ctx := context.WithValue(r.Context(), payloadKey, payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnlyMiddleware memastikan hanya admin yang bisa akses
func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, ok := r.Context().Value(payloadKey).(*auth.PasetoPayload)
		if !ok || payload.Role != "admin" {
			http.Error(w, `{"error": "Admin access required"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetPayloadFromContext mengambil payload dari request context
func GetPayloadFromContext(ctx context.Context) (*auth.PasetoPayload, bool) {
	payload, ok := ctx.Value(payloadKey).(*auth.PasetoPayload)
	return payload, ok
}
