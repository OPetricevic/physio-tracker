package auth

import (
	"encoding/json"
	"net/http"

	"github.com/OPetricevic/physio-tracker/backend/internal/services/auth"
)

type AuthController struct {
	svc auth.Service
}

func NewController(svc auth.Service) *AuthController {
	return &AuthController{svc: svc}
}

type registerRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type loginRequest struct {
	Identifier string `json:"identifier"` // email or username
	Password   string `json:"password"`
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	tok, err := c.svc.Register(r.Context(), &auth.RegisterRequest{
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	})
	if err != nil {
		writeAuthError(w, "register", err)
		return
	}
	writeJSON(w, map[string]interface{}{
		"token":       tok.GetToken(),
		"expires_at":  tok.GetExpiresAt().AsTime(),
		"doctor_uuid": tok.GetDoctorUuid(),
		"token_uuid":  tok.GetUuid(),
	}, http.StatusCreated)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	tok, err := c.svc.Login(r.Context(), req.Identifier, req.Password)
	if err != nil {
		writeAuthError(w, "login", err)
		return
	}
	writeJSON(w, map[string]interface{}{
		"token":       tok.GetToken(),
		"expires_at":  tok.GetExpiresAt().AsTime(),
		"doctor_uuid": tok.GetDoctorUuid(),
		"token_uuid":  tok.GetUuid(),
	}, http.StatusOK)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if err := c.svc.Logout(r.Context(), payload.Token); err != nil {
		writeAuthError(w, "logout", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeAuthError(w http.ResponseWriter, action string, err error) {
	msg := action + ": " + err.Error()
	switch err {
	case auth.ErrInvalidRequest:
		writeJSON(w, map[string]string{"error": "invalid_request", "message": msg}, http.StatusBadRequest)
	case auth.ErrConflict:
		writeJSON(w, map[string]string{"error": "conflict", "message": msg}, http.StatusConflict)
	case auth.ErrUnauthorized:
		writeJSON(w, map[string]string{"error": "unauthorized", "message": msg}, http.StatusUnauthorized)
	default:
		writeJSON(w, map[string]string{"error": "internal_error", "message": msg}, http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, payload interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
