package doctors

import (
	"net/http"

	ctrl "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/doctors"
	"github.com/gorilla/mux"
)

// DoctorHandler wires routes to doctor controller HTTP methods.
type DoctorHandler struct {
	controller *ctrl.DoctorController
}

func NewHandler(controller *ctrl.DoctorController) *DoctorHandler {
	return &DoctorHandler{controller: controller}
}

func (h *DoctorHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/doctors/create", h.controller.CreateDoctor).Methods(http.MethodPost)
	r.HandleFunc("/doctors", h.controller.ListDoctors).Methods(http.MethodGet)
	r.HandleFunc("/doctors/{uuid}", h.controller.UpdateDoctor).Methods(http.MethodPatch)
	r.HandleFunc("/doctors/{uuid}", h.controller.DeleteDoctor).Methods(http.MethodDelete)
}
