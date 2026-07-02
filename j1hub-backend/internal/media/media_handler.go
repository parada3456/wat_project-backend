package media

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/parada3456/wat_project-backend/internal/domain"
	"github.com/parada3456/wat_project-backend/pkg/apperror"
	"github.com/parada3456/wat_project-backend/pkg/uid"
)

type StoragePort interface {
	UploadFile(ctx context.Context, bucket, key string, data io.Reader, contentType string) (url string, err error)
	DeleteFile(ctx context.Context, bucket, key string) error
}

type MediaHandler struct {
	storagePort StoragePort
}

func NewMediaHandler(storagePort StoragePort) *MediaHandler {
	log.Println("debugprint: entering NewMediaHandler")
	return &MediaHandler{
		storagePort: storagePort,
	}
}

func (h *MediaHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*MediaHandler).UploadFile")

	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		apperror.RespondError(w, fmt.Errorf("failed to parse multipart form: %w", domain.ErrInvalidInput))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		apperror.RespondError(w, fmt.Errorf("failed to get file: %w", domain.ErrInvalidInput))
		return
	}
	defer file.Close()

	bucket := r.FormValue("bucket")
	if bucket == "" {
		bucket = r.URL.Query().Get("bucket")
	}
	if bucket == "" {
		bucket = "media"
	}

	key := r.FormValue("key")
	if key == "" {
		key = r.URL.Query().Get("key")
	}
	if key == "" {
		key = uid.New("med_")
		if header != nil {
			parts := strings.Split(header.Filename, ".")
			if len(parts) > 1 {
				ext := parts[len(parts)-1]
				key = fmt.Sprintf("%s.%s", key, ext)
			}
		}
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	url, err := h.storagePort.UploadFile(r.Context(), bucket, key, file, contentType)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"url": url,
		"key": key,
	})
}

func (h *MediaHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*MediaHandler).DeleteFile")

	key := chi.URLParam(r, "key")
	if key == "" {
		apperror.RespondError(w, fmt.Errorf("missing key parameter: %w", domain.ErrInvalidInput))
		return
	}

	bucket := r.URL.Query().Get("bucket")
	if bucket == "" {
		bucket = "media"
	}

	err := h.storagePort.DeleteFile(r.Context(), bucket, key)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
