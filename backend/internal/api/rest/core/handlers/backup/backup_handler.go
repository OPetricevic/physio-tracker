package backup

import (
	"net/http"

	ctrl "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/backup"
	"github.com/gorilla/mux"
)

type Handler struct {
	controller *ctrl.Controller
}

func NewHandler(controller *ctrl.Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/backup", h.controller.Download).Methods(http.MethodGet)
	r.HandleFunc("/backup/restore", h.controller.Restore).Methods(http.MethodPost)
}
