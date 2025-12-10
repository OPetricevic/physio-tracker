package patients

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	se "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/serviceerrors"
	mwauth "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/middleware"
	svc "github.com/OPetricevic/physio-tracker/backend/internal/services/patients"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type PatientController struct {
	svc svc.Service
}

func NewController(svc svc.Service) *PatientController {
	return &PatientController{svc: svc}
}

func parsePositiveInt(val string, def int) int {
	if val == "" {
		return def
	}
	if n, err := strconv.Atoi(val); err == nil && n > 0 {
		return n
	}
	return def
}

// HTTP helpers: controller handles decoding/validation and calls service.
// Handlers only wire routes to these methods.
func (c *PatientController) CreatePatient(w http.ResponseWriter, r *http.Request) {
	doctorUUID, ok := mwauth.GetDoctorUUID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req pb.CreatePatientRequest
	body, _ := io.ReadAll(r.Body)
	if err := jsonpb.Unmarshal(body, &req); err != nil {
		writeJSONError(w, "invalid_request", "create patient: invalid JSON", http.StatusBadRequest)
		return
	}
	req.DoctorUuid = doctorUUID
	p, err := c.svc.Create(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			writeJSONError(w, "invalid_request", "create patient: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			writeJSONError(w, "not_found", "create patient: not found", http.StatusNotFound)
		case errors.Is(err, se.ErrConflict):
			writeJSONError(w, "conflict", "create patient: conflict", http.StatusConflict)
		default:
			writeJSONError(w, "internal_error", "create patient: internal error", http.StatusInternalServerError)
		}
		return
	}
	writeProto(w, p, http.StatusCreated)
}

func (c *PatientController) UpdatePatient(w http.ResponseWriter, r *http.Request) {
	doctorUUID, ok := mwauth.GetDoctorUUID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	patientUUID := vars["uuid"]
	var req pb.UpdatePatientRequest
	body, _ := io.ReadAll(r.Body)
	if err := jsonpb.Unmarshal(body, &req); err != nil {
		writeJSONError(w, "invalid_request", "update patient: invalid JSON", http.StatusBadRequest)
		return
	}
	if patientUUID != "" {
		req.Uuid = patientUUID
	}
	req.DoctorUuid = doctorUUID
	p, err := c.svc.Update(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			writeJSONError(w, "invalid_request", "update patient: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			writeJSONError(w, "not_found", "update patient: not found", http.StatusNotFound)
		case errors.Is(err, se.ErrConflict):
			writeJSONError(w, "conflict", "update patient: conflict", http.StatusConflict)
		default:
			writeJSONError(w, "internal_error", "update patient: internal error", http.StatusInternalServerError)
		}
		return
	}
	writeProto(w, p, http.StatusOK)
}

func (c *PatientController) DeletePatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientUUID := vars["uuid"]
	if err := c.svc.Delete(r.Context(), patientUUID); err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			writeJSONError(w, "invalid_request", "delete patient: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			writeJSONError(w, "not_found", "delete patient: not found", http.StatusNotFound)
		case errors.Is(err, se.ErrConflict):
			writeJSONError(w, "conflict", "delete patient: conflict", http.StatusConflict)
		default:
			writeJSONError(w, "internal_error", "delete patient: internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *PatientController) ListPatients(w http.ResponseWriter, r *http.Request) {
	doctorUUID, ok := mwauth.GetDoctorUUID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	q := r.URL.Query()
	pageSize := parsePositiveInt(q.Get("page_size"), 20)
	currentPage := parsePositiveInt(q.Get("current_page"), 1)
	req := &pb.ListPatientsRequest{
		Query: q.Get("query"),
	}
	list, err := c.svc.List(r.Context(), req, doctorUUID, pageSize, currentPage)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			writeJSONError(w, "invalid_request", "list patients: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			writeJSONError(w, "not_found", "list patients: not found", http.StatusNotFound)
		case errors.Is(err, se.ErrConflict):
			writeJSONError(w, "conflict", "list patients: conflict", http.StatusConflict)
		default:
			writeJSONError(w, "internal_error", "list patients: internal error", http.StatusInternalServerError)
		}
		return
	}
	resp := &pb.ListPatientsResponse{Patients: list}
	writeProto(w, resp, http.StatusOK)
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
