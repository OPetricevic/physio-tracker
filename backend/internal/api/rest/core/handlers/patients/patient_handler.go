package patients

import (
	"net/http"

	ctrl "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/patients"
	"github.com/gorilla/mux"
)

// PatientHandler wires routes to patient controller HTTP methods.
type PatientHandler struct {
	controller *ctrl.PatientController
}

func NewHandler(controller *ctrl.PatientController) *PatientHandler {
	return &PatientHandler{controller: controller}
}

func (h *PatientHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/patients/create", h.controller.CreatePatient).Methods(http.MethodPost)
	r.HandleFunc("/patients", h.controller.ListPatients).Methods(http.MethodGet)
	r.HandleFunc("/patients/{uuid}", h.controller.UpdatePatient).Methods(http.MethodPatch)
	r.HandleFunc("/patients/{uuid}", h.controller.DeletePatient).Methods(http.MethodDelete)
}
