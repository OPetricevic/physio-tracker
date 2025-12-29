package handlers

import (
	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/anamneses"
	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/doctorprofiles"
	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/doctors"
	uploadhandler "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/files"
	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/patients"
	canamneses "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/anamneses"
	cdoctorprofiles "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/doctorprofiles"
	cdoctors "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/doctors"
	cpatients "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/patients"
	dbanamneses "github.com/OPetricevic/physio-tracker/backend/internal/database/anamneses"
	dbdoctorprofiles "github.com/OPetricevic/physio-tracker/backend/internal/database/doctorprofiles"
	dbdoctors "github.com/OPetricevic/physio-tracker/backend/internal/database/doctors"
	dbpatients "github.com/OPetricevic/physio-tracker/backend/internal/database/patients"
	svcanamneses "github.com/OPetricevic/physio-tracker/backend/internal/services/anamneses"
	svcdoctorprofiles "github.com/OPetricevic/physio-tracker/backend/internal/services/doctorprofiles"
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
	NewDoctorProfileModule,
	NewUploadModule,
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
	pRepo := dbpatients.NewPatientsRepository(db)
	profRepo := dbdoctorprofiles.NewRepository(db)
	svc := svcanamneses.NewService(repo, pRepo, profRepo)
	ctrl := canamneses.NewController(svc)
	return &anamnesisModule{handler: anamneses.NewHandler(ctrl)}
}

func (m *anamnesisModule) Register(r *mux.Router) {
	m.handler.RegisterRoutes(r)
}

// Doctor profile module wiring.
type doctorProfileModule struct {
	handler *doctorprofiles.Handler
}

func NewDoctorProfileModule(db *gorm.DB) Module {
	repo := dbdoctorprofiles.NewRepository(db)
	svc := svcdoctorprofiles.NewService(repo)
	ctrl := cdoctorprofiles.NewController(svc)
	return &doctorProfileModule{handler: doctorprofiles.NewHandler(ctrl)}
}

func (m *doctorProfileModule) Register(r *mux.Router) {
	m.handler.RegisterRoutes(r)
}

// Upload module (branding images).
type uploadModule struct {
	handler *uploadhandler.Handler
}

func NewUploadModule(_ *gorm.DB) Module {
	return &uploadModule{handler: uploadhandler.NewHandler()}
}

func (m *uploadModule) Register(r *mux.Router) {
	r.HandleFunc("/files/upload", m.handler.Upload).Methods("POST")
}
