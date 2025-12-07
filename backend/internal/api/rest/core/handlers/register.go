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

// RegisterAll wires all HTTP handlers to the router.
func RegisterAll(r *mux.Router, db *gorm.DB) {
	// Patients
	patientRepo := dbpatients.NewPatientsRepository(db)
	patientSvc := svcpatients.NewService(patientRepo)
	patientController := cpatients.NewController(patientSvc)
	patientHandler := patients.NewHandler(patientController)
	patientHandler.RegisterRoutes(r)

	// Doctors
	doctorRepo := dbdoctors.NewDoctorsRepository(db)
	doctorSvc := svcdoctors.NewService(doctorRepo)
	doctorController := cdoctors.NewController(doctorSvc)
	doctorHandler := doctors.NewHandler(doctorController)
	doctorHandler.RegisterRoutes(r)
}
