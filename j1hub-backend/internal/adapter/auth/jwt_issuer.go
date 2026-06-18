package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/j1hub/backend/internal/port"
)

type jwtIssuer struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTIssuer(cfg *config.Config) port.TokenIssuer {
	return &jwtIssuer{
		secret:        []byte(cfg.JWTSecret),
		accessExpiry:  cfg.JWTExpiry(),
		refreshExpiry: time.Hour * 24 * 7, // 24 hours * 7 days
	}
}

func (i *jwtIssuer) Issue(userID string, isAdmin bool) (*port.TokenPair, error) {
	expiresAt := time.Now().Add(i.accessExpiry)
	accessClaims := jwt.MapClaims{
		"sub":      userID,
		"is_admin": isAdmin,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
		"type":     "access",
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(i.secret)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(i.refreshExpiry).Unix(),
		"iat":  time.Now().Unix(),
		"type": "refresh",
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(i.secret)
	if err != nil {
		return nil, err
	}

	return &port.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (i *jwtIssuer) Verify(tokenString string) (*port.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return i.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != "access" {
			return nil, fmt.Errorf("invalid token type")
		}
		return &port.Claims{
			UserID:  claims["sub"].(string),
			IsAdmin: claims["is_admin"].(bool),
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (i *jwtIssuer) Refresh(refreshTokenString string) (*port.TokenPair, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return i.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != "refresh" {
			return nil, fmt.Errorf("invalid token type")
		}
		userID := claims["sub"].(string)
		// In a real app, check if userID is still valid/not blocked
		// For now, just issue new pair. We don't have isAdmin in refresh claims, 
		// but we could add it or fetch from DB.
		return i.Issue(userID, false) // Default to non-admin for refresh for now
	}

	return nil, fmt.Errorf("invalid refresh token")
}
