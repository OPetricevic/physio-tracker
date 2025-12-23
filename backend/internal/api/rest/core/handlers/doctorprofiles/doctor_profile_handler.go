package doctorprofiles

import (
	"net/http"

	ctrl "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/doctorprofiles"
	"github.com/gorilla/mux"
)

type Handler struct {
	controller *ctrl.Controller
}

func NewHandler(c *ctrl.Controller) *Handler {
	return &Handler{controller: c}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/doctor/profile", h.controller.GetProfile).Methods(http.MethodGet)
	r.HandleFunc("/doctor/profile", h.controller.UpsertProfile).Methods(http.MethodPut, http.MethodPost, http.MethodPatch)
}
