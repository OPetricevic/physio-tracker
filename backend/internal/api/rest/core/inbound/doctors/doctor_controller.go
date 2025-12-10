package doctors

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	se "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/serviceerrors"
	svc "github.com/OPetricevic/physio-tracker/backend/internal/services/doctors"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type DoctorController struct {
	svc svc.Service
}

func NewController(svc svc.Service) *DoctorController {
	return &DoctorController{svc: svc}
}

// HTTP helpers
func (c *DoctorController) CreateDoctor(w http.ResponseWriter, r *http.Request) {
	var req pb.CreateDoctorRequest
	body, _ := io.ReadAll(r.Body)
	if err := jsonpb.Unmarshal(body, &req); err != nil {
		writeJSONError(w, "invalid_request", "create doctor: invalid JSON", http.StatusBadRequest)
		return
	}
	doc, err := c.svc.Create(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			writeJSONError(w, "invalid_request", "create doctor: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrConflict):
			writeJSONError(w, "conflict", "create doctor: conflict", http.StatusConflict)
		default:
			writeJSONError(w, "internal_error", "create doctor: internal error", http.StatusInternalServerError)
		}
		return
	}
	writeProto(w, doc, http.StatusCreated)
}

func (c *DoctorController) UpdateDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorUUID := vars["uuid"]
	var req pb.UpdateDoctorRequest
	body, _ := io.ReadAll(r.Body)
	if err := jsonpb.Unmarshal(body, &req); err != nil {
		writeJSONError(w, "invalid_request", "update doctor: invalid JSON", http.StatusBadRequest)
		return
	}
	if doctorUUID != "" {
		req.Uuid = doctorUUID
	}
	doc, err := c.svc.Update(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			writeJSONError(w, "invalid_request", "update doctor: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			writeJSONError(w, "not_found", "update doctor: not found", http.StatusNotFound)
		case errors.Is(err, se.ErrConflict):
			writeJSONError(w, "conflict", "update doctor: conflict", http.StatusConflict)
		default:
			writeJSONError(w, "internal_error", "update doctor: internal error", http.StatusInternalServerError)
		}
		return
	}
	writeProto(w, doc, http.StatusOK)
}

func (c *DoctorController) DeleteDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorUUID := vars["uuid"]
	if err := c.svc.Delete(r.Context(), doctorUUID); err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			writeJSONError(w, "invalid_request", "delete doctor: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			writeJSONError(w, "not_found", "delete doctor: not found", http.StatusNotFound)
		default:
			writeJSONError(w, "internal_error", "delete doctor: internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// local protojson helper
var jsonpb = &protojson.UnmarshalOptions{DiscardUnknown: true}

func writeProto(w http.ResponseWriter, msg proto.Message, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	b, err := protojson.MarshalOptions{EmitUnpopulated: true, UseEnumNumbers: true}.Marshal(msg)
	if err != nil {
		http.Error(w, "internal_error: failed to encode response", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(b)
}

func writeJSONError(w http.ResponseWriter, code, message string, status int) {
	writeJSON(w, map[string]string{
		"error":   code,
		"message": message,
	}, status)
}

func writeJSON(w http.ResponseWriter, payload interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "internal_error: failed to encode response", http.StatusInternalServerError)
	}
}
