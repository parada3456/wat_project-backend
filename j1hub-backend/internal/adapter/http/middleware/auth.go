package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/pkg/apperror"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const IsAdminKey contextKey = "is_admin"

func Authenticate(issuer port.TokenIssuer) func(http.Handler) http.Handler {
	log.Println("debugprint: entering Authenticate")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			log.Printf("debugprint: Authenticate middleware - incoming Authorization header: %q", authHeader)
			if authHeader == "" {
				log.Println("debugprint: Authenticate middleware - Missing authorization header")
				apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Missing authorization header"})
				return
			}

			// Support both "Bearer <token>" and case-insensitive "bearer <token>"
			var token string
			if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				token = authHeader[7:]
			} else {
				token = authHeader
			}
			log.Printf("debugprint: Authenticate middleware - extracted token: %q", token)

			claims, err := issuer.Verify(token)
			if err != nil {
				log.Printf("debugprint: Authenticate middleware - Verify failed: %v", err)
				apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Invalid token", Err: err})
				return
			}

			log.Printf("debugprint: Authenticate middleware - success, user_id: %s, is_admin: %t", claims.UserID, claims.IsAdmin)
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, IsAdminKey, claims.IsAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAdmin(next http.Handler) http.Handler {
	log.Println("debugprint: entering RequireAdmin")
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
	log.Println("debugprint: entering GetClaims")
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
	log.Println("debugprint: entering ContextWithClaims")
	ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
	return context.WithValue(ctx, IsAdminKey, claims.IsAdmin)
}
