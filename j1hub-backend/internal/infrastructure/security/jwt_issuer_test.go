package security_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/j1hub/backend/internal/infrastructure/security"
	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

func TestJWTIssuer_IssueVerifyAndRefresh(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:      "my_super_secret_signing_key_32_bytes",
		JWTExpiryHours: 2,
	}

	issuer := security.NewJWTIssuer(cfg)
	userID := "usr_123"

	// 1. Issue tokens
	tokens, err := issuer.Issue(userID, true)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.WithinDuration(t, time.Now().Add(2*time.Hour), tokens.ExpiresAt, 5*time.Second)

	// 2. Verify valid access token
	claims, err := issuer.Verify(tokens.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.True(t, claims.IsAdmin)

	// 3. Reject verifying a refresh token as an access token
	_, err = issuer.Verify(tokens.RefreshToken)
	assert.Error(t, err)

	// 4. Refresh tokens
	newTokens, err := issuer.Refresh(tokens.RefreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.RefreshToken)

	// 5. Verify refreshed access token
	newClaims, err := issuer.Verify(newTokens.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, newClaims.UserID)
	assert.False(t, newClaims.IsAdmin) // Default refresh issues non-admin as per implementation

	// 6. Verification failure on invalid token
	_, err = issuer.Verify("invalid_token_string")
	assert.Error(t, err)

	// 7. Reject refreshing an access token as a refresh token
	_, err = issuer.Refresh(tokens.AccessToken)
	assert.Error(t, err)

	// 8. Refresh failure on invalid token
	_, err = issuer.Refresh("invalid_refresh_token_string")
	assert.Error(t, err)

	// 9. Unexpected signing method
	noneToken, _ := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub":      userID,
		"is_admin": true,
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
		"type":     "access",
	}).SignedString(jwt.UnsafeAllowNoneSignatureType)

	_, err = issuer.Verify(noneToken)
	assert.Error(t, err)

	_, err = issuer.Refresh(noneToken)
	assert.Error(t, err)
}
