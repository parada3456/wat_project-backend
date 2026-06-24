package security_test

import (
	"testing"

	"github.com/j1hub/backend/internal/infrastructure/security"
	"github.com/stretchr/testify/assert"
)

func TestArgon2Hasher_HashAndVerify(t *testing.T) {
	hasher := security.NewArgon2Hasher()

	password := "super_secret_password"
	hash, err := hasher.Hash(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Verify correct password
	assert.True(t, hasher.Verify(password, hash))

	// Verify incorrect password
	assert.False(t, hasher.Verify("wrong_password", hash))

	// Verify malformed hash string format
	assert.False(t, hasher.Verify(password, "$invalid$hash$format"))
}
