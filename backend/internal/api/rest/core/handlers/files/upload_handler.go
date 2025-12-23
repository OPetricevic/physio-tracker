package files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	mwauth "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/middleware"
)

// Handler handles simple image uploads (logo/signature). Stores under uploads/branding and returns a URL.
type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Upload expects multipart/form-data with field "file"; returns JSON { "url": "/static/branding/<name>" }.
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if _, ok := mwauth.GetDoctorUUID(r.Context()); !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// limit to 5MB
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	switch ext {
	case ".png", ".jpg", ".jpeg":
	default:
		http.Error(w, "unsupported file type", http.StatusBadRequest)
		return
	}

	dstDir := filepath.Join("uploads", "branding")
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		http.Error(w, "unable to store file", http.StatusInternalServerError)
		return
	}
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dstPath := filepath.Join(dstDir, name)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "unable to store file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "unable to store file", http.StatusInternalServerError)
		return
	}

	url := "/static/branding/" + name
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"url":"%s"}`, url)))
}
