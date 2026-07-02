package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	port "github.com/parada3456/wat_project-backend/internal/auth/port"
	"golang.org/x/crypto/argon2"
)

type argon2Hasher struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

func NewArgon2Hasher() port.PasswordHasher {
	log.Println("debugprint: entering NewArgon2Hasher")
	return &argon2Hasher{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
		saltLen: 16,
	}
}

func (h *argon2Hasher) Hash(plain string) (string, error) {
	log.Println("debugprint: entering (*argon2Hasher).Hash")
	salt := make([]byte, h.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(plain), salt, h.time, h.memory, h.threads, h.keyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, h.memory, h.time, h.threads, b64Salt, b64Hash)
	return encoded, nil
}

func (h *argon2Hasher) Verify(plain, hash string) bool {
	log.Println("debugprint: entering (*argon2Hasher).Verify")
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false
	}

	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	keyLen := uint32(len(decodedHash))
	comparisonHash := argon2.IDKey([]byte(plain), salt, time, memory, threads, keyLen)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1
}
