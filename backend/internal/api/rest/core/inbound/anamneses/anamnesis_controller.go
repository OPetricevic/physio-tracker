package anamneses

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	common "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/common"
	mwauth "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/middleware"
	se "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/serviceerrors"
	svc "github.com/OPetricevic/physio-tracker/backend/internal/services/anamneses"
	"github.com/gorilla/mux"
)

type Controller struct {
	svc svc.Service
}

func NewController(s svc.Service) *Controller {
	return &Controller{svc: s}
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

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	doctorUUID, ok := mwauth.GetDoctorUUID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req pb.CreateAnamnesisRequest
	body, _ := io.ReadAll(r.Body)
	if err := common.JSONPB.Unmarshal(body, &req); err != nil {
		common.WriteJSONError(w, "invalid_request", "create anamneza: invalid JSON", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	if patientUUID, ok := vars["patient_uuid"]; ok && patientUUID != "" {
		req.PatientUuid = patientUUID
	}
	if err := common.ValidateProto(&req); err != nil {
		common.WriteJSONError(w, "invalid_request", "create anamneza: "+err.Error(), http.StatusBadRequest)
		return
	}
	anm, err := c.svc.Create(r.Context(), doctorUUID, &req)
	if err != nil {
		switch {
		case isSvcErr(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", err.Error(), http.StatusBadRequest)
		case isSvcErr(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", err.Error(), http.StatusNotFound)
		case isSvcErr(err, se.ErrConflict):
			common.WriteJSONError(w, "conflict", err.Error(), http.StatusConflict)
		default:
			common.WriteJSONError(w, "internal_error", err.Error(), http.StatusInternalServerError)
		}
		return
	}
	common.WriteProto(w, anm, http.StatusCreated)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	doctorUUID, ok := mwauth.GetDoctorUUID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	anamnesisUUID := vars["uuid"]
	patientUUID := vars["patient_uuid"]

	var req pb.UpdateAnamnesisRequest
	body, _ := io.ReadAll(r.Body)
	if err := common.JSONPB.Unmarshal(body, &req); err != nil {
		common.WriteJSONError(w, "invalid_request", "update anamneza: invalid JSON", http.StatusBadRequest)
		return
	}
	if patientUUID != "" {
		req.PatientUuid = patientUUID
	}
	if anamnesisUUID != "" {
		req.Uuid = anamnesisUUID
	}
	if err := common.ValidateProto(&req); err != nil {
		common.WriteJSONError(w, "invalid_request", "update anamneza: "+err.Error(), http.StatusBadRequest)
		return
	}
	anm, err := c.svc.Update(r.Context(), doctorUUID, &req)
	if err != nil {
		switch {
		case isSvcErr(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", err.Error(), http.StatusBadRequest)
		case isSvcErr(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", err.Error(), http.StatusNotFound)
		case isSvcErr(err, se.ErrConflict):
			common.WriteJSONError(w, "conflict", err.Error(), http.StatusConflict)
		default:
			common.WriteJSONError(w, "internal_error", err.Error(), http.StatusInternalServerError)
		}
		return
	}
	common.WriteProto(w, anm, http.StatusOK)
}

func (c *Controller) List(w http.ResponseWriter, r *http.Request) {
	doctorUUID, ok := mwauth.GetDoctorUUID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	patientUUID := vars["patient_uuid"]
	q := r.URL.Query()
	pageSize := parsePositiveInt(q.Get("page_size"), 5)
	currentPage := parsePositiveInt(q.Get("current_page"), 1)
	query := q.Get("query")

	list, err := c.svc.List(r.Context(), doctorUUID, patientUUID, query, pageSize, currentPage)
	if err != nil {
		switch {
		case isSvcErr(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", err.Error(), http.StatusBadRequest)
		case isSvcErr(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", err.Error(), http.StatusNotFound)
		default:
			common.WriteJSONError(w, "internal_error", err.Error(), http.StatusInternalServerError)
		}
		return
	}
	resp := &pb.ListAnamnesesResponse{Anamneses: list}
	common.WriteProto(w, resp, http.StatusOK)
}

// GeneratePDF builds a PDF for a specific anamnesis.
func (c *Controller) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	doctorUUID, ok := mwauth.GetDoctorUUID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	patientUUID := vars["patient_uuid"]
	anamnesisUUID := vars["uuid"]

	var req pb.GenerateAnamnesisPdfRequest
	body, _ := io.ReadAll(r.Body)
	if len(body) > 0 {
		_ = common.JSONPB.Unmarshal(body, &req)
	}

	onlyCurrent := r.URL.Query().Get("only_current") == "true"

	bytes, err := c.svc.GeneratePDF(r.Context(), doctorUUID, patientUUID, anamnesisUUID, req.IncludeVisitUuids, onlyCurrent)
	if err != nil {
		// Log full error for debugging 500s.
		log.Printf("generate pdf failed: %v", err)
		switch {
		case isSvcErr(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", err.Error(), http.StatusBadRequest)
		case isSvcErr(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", err.Error(), http.StatusNotFound)
		default:
			common.WriteJSONError(w, "internal_error", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"anamnesis.pdf\"")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bytes)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	doctorUUID, ok := mwauth.GetDoctorUUID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	anamnesisUUID := vars["uuid"]
	if err := c.svc.Delete(r.Context(), doctorUUID, anamnesisUUID); err != nil {
		switch {
		case isSvcErr(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", err.Error(), http.StatusBadRequest)
		case isSvcErr(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", err.Error(), http.StatusNotFound)
		default:
			common.WriteJSONError(w, "internal_error", err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func isSvcErr(err error, target error) bool {
	return errors.Is(err, target)
}
