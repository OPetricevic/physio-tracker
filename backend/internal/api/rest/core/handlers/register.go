package handlers

import (
	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/doctors"
	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/patients"
	cdoctors "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/doctors"
	cpatients "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/patients"
	dbdoctors "github.com/OPetricevic/physio-tracker/backend/internal/database/doctors"
	dbpatients "github.com/OPetricevic/physio-tracker/backend/internal/database/patients"
	svcdoctors "github.com/OPetricevic/physio-tracker/backend/internal/services/doctors"
	svcpatients "github.com/OPetricevic/physio-tracker/backend/internal/services/patients"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// module registers routes for a feature area.
type module interface {
	Register(r *mux.Router)
}

func newPatientModule(db *gorm.DB) module {
	repo := dbpatients.NewPatientsRepository(db)
	svc := svcpatients.NewService(repo)
	ctrl := cpatients.NewController(svc)
	h := patients.NewHandler(ctrl)
	return h
}

func newDoctorModule(db *gorm.DB) module {
	repo := dbdoctors.NewDoctorsRepository(db)
	svc := svcdoctors.NewService(repo)
	ctrl := cdoctors.NewController(svc)
	h := doctors.NewHandler(ctrl)
	return h
}

// RegisterAll wires all HTTP handlers to the router.
func RegisterAll(r *mux.Router, db *gorm.DB) {
	modules := []module{
		newPatientModule(db),
		newDoctorModule(db),
	}
	for _, m := range modules {
		m.Register(r)
	}
}
