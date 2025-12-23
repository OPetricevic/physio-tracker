package handlers

import (
	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/doctors"
	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/anamneses"
	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/patients"
	canamneses "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/anamneses"
	cdoctors "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/doctors"
	cpatients "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/patients"
	dbanamneses "github.com/OPetricevic/physio-tracker/backend/internal/database/anamneses"
	dbdoctors "github.com/OPetricevic/physio-tracker/backend/internal/database/doctors"
	dbpatients "github.com/OPetricevic/physio-tracker/backend/internal/database/patients"
	svcanamneses "github.com/OPetricevic/physio-tracker/backend/internal/services/anamneses"
	svcdoctors "github.com/OPetricevic/physio-tracker/backend/internal/services/doctors"
	svcpatients "github.com/OPetricevic/physio-tracker/backend/internal/services/patients"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Module registers HTTP routes for a feature area.
type Module interface {
	Register(r *mux.Router)
}

// moduleBuilders is the registry: add new modules here to expose new APIs.
var moduleBuilders = []func(*gorm.DB) Module{
	NewPatientModule,
	NewDoctorModule,
	NewAnamnesisModule,
}

// Patient module wiring (repo -> service -> controller -> handler).
type patientModule struct {
	handler *patients.PatientHandler
}

func NewPatientModule(db *gorm.DB) Module {
	repo := dbpatients.NewPatientsRepository(db)
	svc := svcpatients.NewService(repo)
	ctrl := cpatients.NewController(svc)
	return &patientModule{handler: patients.NewHandler(ctrl)}
}

func (m *patientModule) Register(r *mux.Router) {
	m.handler.RegisterRoutes(r)
}

// Doctor module wiring.
type doctorModule struct {
	handler *doctors.DoctorHandler
}

func NewDoctorModule(db *gorm.DB) Module {
	repo := dbdoctors.NewDoctorsRepository(db)
	svc := svcdoctors.NewService(repo)
	ctrl := cdoctors.NewController(svc)
	return &doctorModule{handler: doctors.NewHandler(ctrl)}
}

func (m *doctorModule) Register(r *mux.Router) {
	m.handler.RegisterRoutes(r)
}

// Anamnesis module wiring.
type anamnesisModule struct {
	handler *anamneses.Handler
}

func NewAnamnesisModule(db *gorm.DB) Module {
	repo := dbanamneses.NewRepository(db)
	svc := svcanamneses.NewService(repo)
	ctrl := canamneses.NewController(svc)
	return &anamnesisModule{handler: anamneses.NewHandler(ctrl)}
}

func (m *anamnesisModule) Register(r *mux.Router) {
	m.handler.RegisterRoutes(r)
}
