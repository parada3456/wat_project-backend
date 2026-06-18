package port

import (
	"context"
	"io"
	"time"
)

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Verify(plain, hash string) bool
}

type Claims struct {
	UserID  string
	IsAdmin bool
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type TokenIssuer interface {
	Issue(userID string, isAdmin bool) (*TokenPair, error)
	Verify(token string) (*Claims, error)
	Refresh(refreshToken string) (*TokenPair, error)
}

type StoragePort interface {
	UploadFile(ctx context.Context, bucket, key string, data io.Reader, contentType string) (url string, err error)
}

type NotifierPort interface {
	Send(ctx context.Context, userID, title, body string) error
}
