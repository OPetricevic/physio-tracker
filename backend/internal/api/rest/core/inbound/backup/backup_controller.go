package backup

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/common"
	"github.com/OPetricevic/physio-tracker/backend/internal/services/backup"
)

type Controller struct {
	svc backup.Service
}

func NewController(svc backup.Service) *Controller {
	return &Controller{svc: svc}
}

func (c *Controller) Download(w http.ResponseWriter, r *http.Request) {
	name, path, err := c.svc.CreateBackup(r.Context())
	if err != nil {
		common.WriteJSONError(w, "internal_error", err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(path)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	http.ServeFile(w, r, path)
}

func (c *Controller) Restore(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		common.WriteJSONError(w, "invalid_request", "neispravan upload", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		common.WriteJSONError(w, "invalid_request", "nedostaje datoteka", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tmp, err := os.CreateTemp("", "physio-restore-*.dump")
	if err != nil {
		common.WriteJSONError(w, "internal_error", "nije moguće spremiti privremenu datoteku", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmp.Name())
	if _, err := io.Copy(tmp, file); err != nil {
		common.WriteJSONError(w, "internal_error", "neuspješno spremanje datoteke", http.StatusInternalServerError)
		return
	}
	_ = tmp.Close()

	if err := c.svc.RestoreBackup(r.Context(), tmp.Name()); err != nil {
		common.WriteJSONError(w, "internal_error", err.Error(), http.StatusInternalServerError)
		return
	}
	common.WriteJSON(w, map[string]string{
		"message": "Sigurnosna kopija je vraćena.",
		"file":    filepath.Base(header.Filename),
	}, http.StatusOK)
}
