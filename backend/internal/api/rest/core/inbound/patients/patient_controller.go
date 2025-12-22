package patients

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	"google.golang.org/protobuf/encoding/protojson"

	common "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/common"
	mwauth "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/middleware"
	se "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/serviceerrors"
	svc "github.com/OPetricevic/physio-tracker/backend/internal/services/patients"
	"github.com/gorilla/mux"
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
		common.WriteJSONError(w, "invalid_request", "create patient: invalid JSON", http.StatusBadRequest)
		return
	}
	req.DoctorUuid = doctorUUID // enforce ownership
	if err := common.ValidateProto(&req); err != nil {
		common.WriteJSONError(w, "invalid_request", "create patient: "+err.Error(), http.StatusBadRequest)
		return
	}
	p, err := c.svc.Create(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", "create patient: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", "create patient: not found", http.StatusNotFound)
		case errors.Is(err, se.ErrConflict):
			common.WriteJSONError(w, "conflict", "create patient: conflict", http.StatusConflict)
		default:
			common.WriteJSONError(w, "internal_error", "create patient: internal error", http.StatusInternalServerError)
		}
		return
	}
	common.WriteProto(w, p, http.StatusCreated)
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
		common.WriteJSONError(w, "invalid_request", "update patient: invalid JSON", http.StatusBadRequest)
		return
	}
	req.DoctorUuid = doctorUUID // enforce ownership
	if err := common.ValidateProto(&req); err != nil {
		common.WriteJSONError(w, "invalid_request", "update patient: "+err.Error(), http.StatusBadRequest)
		return
	}
	if patientUUID != "" {
		req.Uuid = patientUUID
	}
	p, err := c.svc.Update(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", "update patient: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", "update patient: not found", http.StatusNotFound)
		case errors.Is(err, se.ErrConflict):
			common.WriteJSONError(w, "conflict", "update patient: conflict", http.StatusConflict)
		default:
			common.WriteJSONError(w, "internal_error", "update patient: internal error", http.StatusInternalServerError)
		}
		return
	}
	common.WriteProto(w, p, http.StatusOK)
}

func (c *PatientController) DeletePatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientUUID := vars["uuid"]
	if err := c.svc.Delete(r.Context(), patientUUID); err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", "delete patient: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", "delete patient: not found", http.StatusNotFound)
		case errors.Is(err, se.ErrConflict):
			common.WriteJSONError(w, "conflict", "delete patient: conflict", http.StatusConflict)
		default:
			common.WriteJSONError(w, "internal_error", "delete patient: internal error", http.StatusInternalServerError)
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
			common.WriteJSONError(w, "invalid_request", "list patients: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", "list patients: not found", http.StatusNotFound)
		case errors.Is(err, se.ErrConflict):
			common.WriteJSONError(w, "conflict", "list patients: conflict", http.StatusConflict)
		default:
			common.WriteJSONError(w, "internal_error", "list patients: internal error", http.StatusInternalServerError)
		}
		return
	}
	resp := &pb.ListPatientsResponse{Patients: list}
	common.WriteProto(w, resp, http.StatusOK)
}

var jsonpb = &protojson.UnmarshalOptions{DiscardUnknown: true}
