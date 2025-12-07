package main

import (
	"log"
	"net/http"
	"os"

	hdoctors "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/doctors"
	hpatients "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/patients"
	cdoctors "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/doctors"
	cpatients "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/patients"
	dbdoctors "github.com/OPetricevic/physio-tracker/backend/internal/database/doctors"
	dbpatients "github.com/OPetricevic/physio-tracker/backend/internal/database/patients"
	svcdoctors "github.com/OPetricevic/physio-tracker/backend/internal/services/doctors"
	svcpatients "github.com/OPetricevic/physio-tracker/backend/internal/services/patients"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	addr := envOrDefault("PORT", "3600")

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	dsn := envOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/physio?sslmode=disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// Wire patients handler (Postgres via GORM).
	patientRepo := dbpatients.NewPatientsRepository(db)
	patientSvc := svcpatients.NewService(patientRepo)
	patientController := cpatients.NewController(patientSvc)
	patientHandler := hpatients.NewHandler(patientController)
	patientHandler.RegisterRoutes(r)

	// Wire doctors handler
	doctorRepo := dbdoctors.NewDoctorsRepository(db)
	doctorSvc := svcdoctors.NewService(doctorRepo)
	doctorController := cdoctors.NewController(doctorSvc)
	doctorHandler := hdoctors.NewHandler(doctorController)
	doctorHandler.RegisterRoutes(r)

	server := &http.Server{
		Addr:    ":" + addr,
		Handler: r,
	}

	log.Printf("listening on :%s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
