package auth

import (
	"net/http"

	ctrl "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/auth"
	"github.com/gorilla/mux"
)

// AuthHandler wires auth routes.
type AuthHandler struct {
	controller *ctrl.AuthController
}

func NewHandler(controller *ctrl.AuthController) *AuthHandler {
	return &AuthHandler{controller: controller}
}

func (h *AuthHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/auth/register", h.controller.Register).Methods(http.MethodPost)
	r.HandleFunc("/auth/login", h.controller.Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/logout", h.controller.Logout).Methods(http.MethodPost)
}
