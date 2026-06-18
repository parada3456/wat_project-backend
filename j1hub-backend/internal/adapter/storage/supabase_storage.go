package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/j1hub/backend/internal/port"
)

type supabaseStorage struct {
	url string
	key string
}

func NewSupabaseStorage(cfg *config.Config) port.StoragePort {
	return &supabaseStorage{
		url: cfg.SupabaseURL,
		key: cfg.SupabaseServiceKey,
	}
}

func (s *supabaseStorage) UploadFile(ctx context.Context, bucket, key string, data io.Reader, contentType string) (string, error) {
	// Simple implementation using Supabase Storage REST API
	// URL: {url}/storage/v1/object/{bucket}/{key}
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.url, bucket, key)

	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, data)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.key)
	req.Header.Set("apikey", s.key)
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("storage upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Public URL: {url}/storage/v1/object/public/{bucket}/{key}
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.url, bucket, key), nil
}
