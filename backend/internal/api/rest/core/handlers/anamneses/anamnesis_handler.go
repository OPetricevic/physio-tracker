package anamneses

import (
	"net/http"

	ctrl "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/anamneses"
	"github.com/gorilla/mux"
)

type Handler struct {
	controller *ctrl.Controller
}

func NewHandler(controller *ctrl.Controller) *Handler {
	return &Handler{controller: controller}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	// All routes are protected (doctor auth) via parent router middleware.
	r.HandleFunc("/patients/{patient_uuid}/anamneses", h.controller.List).Methods(http.MethodGet)
	r.HandleFunc("/patients/{patient_uuid}/anamneses", h.controller.Create).Methods(http.MethodPost)
	r.HandleFunc("/patients/{patient_uuid}/anamneses/{uuid}", h.controller.Update).Methods(http.MethodPatch)
	r.HandleFunc("/patients/{patient_uuid}/anamneses/{uuid}", h.controller.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/patients/{patient_uuid}/anamneses/{uuid}/pdf", h.controller.GeneratePDF).Methods(http.MethodPost)
}
