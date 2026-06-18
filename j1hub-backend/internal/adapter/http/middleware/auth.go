package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/pkg/apperror"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const IsAdminKey contextKey = "is_admin"

func Authenticate(issuer port.TokenIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Missing authorization header"})
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := issuer.Verify(token)
			if err != nil {
				apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Invalid token", Err: err})
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, IsAdminKey, claims.IsAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value(IsAdminKey).(bool)
		if !ok || !isAdmin {
			apperror.RespondError(w, &apperror.AppError{Code: http.StatusForbidden, Message: "Admin access required"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetClaims(ctx context.Context) *port.Claims {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return nil
	}
	isAdmin, _ := ctx.Value(IsAdminKey).(bool)
	return &port.Claims{
		UserID:  userID,
		IsAdmin: isAdmin,
	}
}

func ContextWithClaims(ctx context.Context, claims *port.Claims) context.Context {
	ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
	return context.WithValue(ctx, IsAdminKey, claims.IsAdmin)
}
