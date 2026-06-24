package media_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/j1hub/backend/internal/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStoragePort struct {
	mock.Mock
}

func (m *MockStoragePort) UploadFile(ctx context.Context, bucket, key string, data io.Reader, contentType string) (string, error) {
	args := m.Called(ctx, bucket, key, data, contentType)
	return args.String(0), args.Error(1)
}

func (m *MockStoragePort) DeleteFile(ctx context.Context, bucket, key string) error {
	args := m.Called(ctx, bucket, key)
	return args.Error(0)
}

func TestMediaHandler_UploadFile_Success(t *testing.T) {
	mockStorage := new(MockStoragePort)
	handler := media.NewMediaHandler(mockStorage)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.png")
	assert.NoError(t, err)
	_, err = part.Write([]byte("image-data"))
	assert.NoError(t, err)
	err = writer.WriteField("bucket", "custom-bucket")
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)

	mockStorage.On("UploadFile", mock.Anything, "custom-bucket", mock.Anything, mock.Anything, "application/octet-stream").
		Return("https://supabase.co/storage/v1/object/public/custom-bucket/med_test.png", nil)

	req := httptest.NewRequest("POST", "/api/v1/media/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler.UploadFile(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "https://supabase.co/storage/v1/object/public/custom-bucket/med_test.png", resp["url"])
	assert.NotEmpty(t, resp["key"])
	mockStorage.AssertExpectations(t)
}

func TestMediaHandler_UploadFile_MissingFile(t *testing.T) {
	mockStorage := new(MockStoragePort)
	handler := media.NewMediaHandler(mockStorage)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err := writer.Close()
	assert.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/media/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler.UploadFile(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestMediaHandler_DeleteFile_Success(t *testing.T) {
	mockStorage := new(MockStoragePort)
	handler := media.NewMediaHandler(mockStorage)

	mockStorage.On("DeleteFile", mock.Anything, "media", "med_123.png").Return(nil)

	req := httptest.NewRequest("DELETE", "/api/v1/media/med_123.png", nil)
	
	// Add chi route context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("key", "med_123.png")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.DeleteFile(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockStorage.AssertExpectations(t)
}

func TestMediaHandler_DeleteFile_StorageError(t *testing.T) {
	mockStorage := new(MockStoragePort)
	handler := media.NewMediaHandler(mockStorage)

	mockStorage.On("DeleteFile", mock.Anything, "custom", "med_123.png").Return(errors.New("storage error"))

	req := httptest.NewRequest("DELETE", "/api/v1/media/med_123.png?bucket=custom", nil)
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("key", "med_123.png")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.DeleteFile(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockStorage.AssertExpectations(t)
}
