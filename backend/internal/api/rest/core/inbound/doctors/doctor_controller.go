package doctors

import (
	"errors"
	"io"
	"net/http"

	pb "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	common "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/common"
	mwauth "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/middleware"
	se "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/serviceerrors"
	svc "github.com/OPetricevic/physio-tracker/backend/internal/services/doctors"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/encoding/protojson"
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
		common.WriteJSONError(w, "invalid_request", "create doctor: invalid JSON", http.StatusBadRequest)
		return
	}
	doc, err := c.svc.Create(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", "create doctor: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrConflict):
			common.WriteJSONError(w, "conflict", "create doctor: conflict", http.StatusConflict)
		default:
			common.WriteJSONError(w, "internal_error", "create doctor: internal error", http.StatusInternalServerError)
		}
		return
	}
	common.WriteProto(w, doc, http.StatusCreated)
}

func (c *DoctorController) UpdateDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorUUID := vars["uuid"]
	var req pb.UpdateDoctorRequest
	body, _ := io.ReadAll(r.Body)
	if err := jsonpb.Unmarshal(body, &req); err != nil {
		common.WriteJSONError(w, "invalid_request", "update doctor: invalid JSON", http.StatusBadRequest)
		return
	}
	if doctorUUID != "" {
		req.Uuid = doctorUUID
	}
	doc, err := c.svc.Update(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", "update doctor: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", "update doctor: not found", http.StatusNotFound)
		case errors.Is(err, se.ErrConflict):
			common.WriteJSONError(w, "conflict", "update doctor: conflict", http.StatusConflict)
		default:
			common.WriteJSONError(w, "internal_error", "update doctor: internal error", http.StatusInternalServerError)
		}
		return
	}
	common.WriteProto(w, doc, http.StatusOK)
}

func (c *DoctorController) GetMe(w http.ResponseWriter, r *http.Request) {
	doctorUUID, ok := mwauth.GetDoctorUUID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	doc, err := c.svc.Get(r.Context(), doctorUUID)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", "get doctor: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", "get doctor: not found", http.StatusNotFound)
		default:
			common.WriteJSONError(w, "internal_error", "get doctor: internal error", http.StatusInternalServerError)
		}
		return
	}
	common.WriteProto(w, doc, http.StatusOK)
}

func (c *DoctorController) UpdateMe(w http.ResponseWriter, r *http.Request) {
	doctorUUID, ok := mwauth.GetDoctorUUID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req pb.UpdateDoctorRequest
	body, _ := io.ReadAll(r.Body)
	if err := jsonpb.Unmarshal(body, &req); err != nil {
		common.WriteJSONError(w, "invalid_request", "update doctor: invalid JSON", http.StatusBadRequest)
		return
	}
	req.Uuid = doctorUUID
	doc, err := c.svc.Update(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", "update doctor: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", "update doctor: not found", http.StatusNotFound)
		case errors.Is(err, se.ErrConflict):
			common.WriteJSONError(w, "conflict", "update doctor: conflict", http.StatusConflict)
		default:
			common.WriteJSONError(w, "internal_error", "update doctor: internal error", http.StatusInternalServerError)
		}
		return
	}
	common.WriteProto(w, doc, http.StatusOK)
}

func (c *DoctorController) DeleteDoctor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	doctorUUID := vars["uuid"]
	if err := c.svc.Delete(r.Context(), doctorUUID); err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", "delete doctor: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", "delete doctor: not found", http.StatusNotFound)
		default:
			common.WriteJSONError(w, "internal_error", "delete doctor: internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// local protojson helper (matches patients controller style)
var jsonpb = &protojson.UnmarshalOptions{DiscardUnknown: true}
