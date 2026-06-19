package auth

import (
	"fmt"
	"log"
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
	log.Println("debugprint: entering NewJWTIssuer")
	return &jwtIssuer{
		secret:        []byte(cfg.JWTSecret),
		accessExpiry:  cfg.JWTExpiry(),
		refreshExpiry: time.Hour * 24 * 7, // 24 hours * 7 days
	}
}

func (i *jwtIssuer) Issue(userID string, isAdmin bool) (*port.TokenPair, error) {
	log.Println("debugprint: entering (*jwtIssuer).Issue")
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
	log.Printf("debugprint: entering (*jwtIssuer).Verify with token: %q", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return i.secret, nil
	})

	if err != nil {
		log.Printf("debugprint: (*jwtIssuer).Verify Parse error: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Printf("debugprint: (*jwtIssuer).Verify claims parsed successfully: %+v", claims)
		if claims["type"] != "access" {
			log.Printf("debugprint: (*jwtIssuer).Verify invalid token type: expected 'access', got %q", claims["type"])
			return nil, fmt.Errorf("invalid token type")
		}
		
		userID, _ := claims["sub"].(string)
		isAdmin, _ := claims["is_admin"].(bool)
		return &port.Claims{
			UserID:  userID,
			IsAdmin: isAdmin,
		}, nil
	}

	log.Println("debugprint: (*jwtIssuer).Verify - token parsed but claims are invalid or token not Valid")
	return nil, fmt.Errorf("invalid token")
}


func (i *jwtIssuer) Refresh(refreshTokenString string) (*port.TokenPair, error) {
	log.Println("debugprint: entering (*jwtIssuer).Refresh")
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
