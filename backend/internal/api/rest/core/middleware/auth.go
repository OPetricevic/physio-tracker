package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	atokens "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/outbound/auth"
)

type contextKey string

const doctorUUIDKey contextKey = "doctor_uuid"

// GetDoctorUUID retrieves the authenticated doctor UUID from context.
func GetDoctorUUID(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(doctorUUIDKey).(string)
	return val, ok && strings.TrimSpace(val) != ""
}

// AuthMiddleware resolves Bearer token -> doctor_uuid and injects into context.
func AuthMiddleware(tokens atokens.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			if !strings.HasPrefix(strings.ToLower(authz), "bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			token := strings.TrimSpace(authz[7:])
			if token == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			t, err := tokens.Get(r.Context(), token)
			if err != nil || t.ExpiresAt.AsTime().Before(time.Now()) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), doctorUUIDKey, t.GetDoctorUuid())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
