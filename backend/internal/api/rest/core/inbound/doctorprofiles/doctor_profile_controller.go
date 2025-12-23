package doctorprofiles

import (
	"errors"
	"io"
	"net/http"

	pt "github.com/OPetricevic/physio-tracker/backend/golang/patients"
	common "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/common"
	mwauth "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/middleware"
	se "github.com/OPetricevic/physio-tracker/backend/internal/commonerrors/serviceerrors"
	svc "github.com/OPetricevic/physio-tracker/backend/internal/services/doctorprofiles"
)

type Controller struct {
	svc svc.Service
}

func NewController(svc svc.Service) *Controller {
	return &Controller{svc: svc}
}

func (c *Controller) GetProfile(w http.ResponseWriter, r *http.Request) {
	doctorUUID, ok := mwauth.GetDoctorUUID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	profile, err := c.svc.Get(r.Context(), doctorUUID)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", "get doctor profile: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", "get doctor profile: not found", http.StatusNotFound)
		default:
			common.WriteJSONError(w, "internal_error", "get doctor profile: internal error", http.StatusInternalServerError)
		}
		return
	}
	resp := &pt.GetDoctorProfileResponse{Profile: profile}
	common.WriteProto(w, resp, http.StatusOK)
}

func (c *Controller) UpsertProfile(w http.ResponseWriter, r *http.Request) {
	doctorUUID, ok := mwauth.GetDoctorUUID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req pt.UpsertDoctorProfileRequest
	body, _ := io.ReadAll(r.Body)
	if err := common.JSONPB.Unmarshal(body, &req); err != nil {
		common.WriteJSONError(w, "invalid_request", "upsert doctor profile: invalid JSON", http.StatusBadRequest)
		return
	}
	// enforce ownership
	if req.Profile == nil {
		req.Profile = &pt.DoctorProfile{}
	}
	req.Profile.DoctorUuid = doctorUUID
	if err := common.ValidateProto(&req); err != nil {
		common.WriteJSONError(w, "invalid_request", "upsert doctor profile: "+err.Error(), http.StatusBadRequest)
		return
	}
	profile, err := c.svc.Upsert(r.Context(), doctorUUID, &req)
	if err != nil {
		switch {
		case errors.Is(err, se.ErrInvalidRequest):
			common.WriteJSONError(w, "invalid_request", "upsert doctor profile: invalid request", http.StatusBadRequest)
		case errors.Is(err, se.ErrNotFound):
			common.WriteJSONError(w, "not_found", "upsert doctor profile: not found", http.StatusNotFound)
		default:
			common.WriteJSONError(w, "internal_error", "upsert doctor profile: internal error", http.StatusInternalServerError)
		}
		return
	}
	resp := &pt.GetDoctorProfileResponse{Profile: profile}
	common.WriteProto(w, resp, http.StatusOK)
}
