package storage_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/j1hub/backend/internal/adapter/storage"
	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

func TestSupabaseStorage_UploadFile(t *testing.T) {
	// Success path
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "Bearer test_key", r.Header.Get("Authorization"))
		assert.Equal(t, "test_key", r.Header.Get("apikey"))
		assert.Equal(t, "image/jpeg", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Key": "bucket/test.jpg"}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		SupabaseURL:        server.URL,
		SupabaseServiceKey: "test_key",
	}

	storagePort := storage.NewSupabaseStorage(cfg)
	url, err := storagePort.UploadFile(context.Background(), "bucket", "test.jpg", strings.NewReader("data"), "image/jpeg")
	assert.NoError(t, err)
	assert.Contains(t, url, "/storage/v1/object/public/bucket/test.jpg")

	// Error status path
	errServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}))
	defer errServer.Close()

	cfgErr := &config.Config{
		SupabaseURL:        errServer.URL,
		SupabaseServiceKey: "test_key",
	}
	storagePortErr := storage.NewSupabaseStorage(cfgErr)
	url, err = storagePortErr.UploadFile(context.Background(), "bucket", "test.jpg", strings.NewReader("data"), "image/jpeg")
	assert.Error(t, err)
	assert.Empty(t, url)
}
