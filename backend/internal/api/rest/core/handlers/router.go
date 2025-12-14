package handlers

import (
	"net/http"

	authh "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/auth"
	authctrl "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/auth"
	mwauth "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/middleware"
	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/middleware"
	dbauth "github.com/OPetricevic/physio-tracker/backend/internal/database/auth"
	dbcreds "github.com/OPetricevic/physio-tracker/backend/internal/database/doctors"
	dbdoctors "github.com/OPetricevic/physio-tracker/backend/internal/database/doctors"
	authsvc "github.com/OPetricevic/physio-tracker/backend/internal/services/auth"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// BuildRouter constructs the HTTP router with health, auth (public), and protected routes.
func BuildRouter(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()
	// order matters: recover first, then logging.
	r.Use(mwauth.RecoverMiddleware)
	r.Use(mwauth.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)

	// Health
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	// Auth wiring (public)
	tokenRepo := dbauth.NewTokensRepository(db)
	credRepo := dbcreds.NewCredentialsRepository(db)
	doctorRepo := dbdoctors.NewDoctorsRepository(db)
	authService := authsvc.NewService(doctorRepo, credRepo, tokenRepo)
	authController := authctrl.NewController(authService)
	authHandler := authh.NewHandler(authController)
	authHandler.RegisterRoutes(r)

	// Protected routes: everything else goes under a subrouter with auth middleware.
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(mwauth.AuthMiddleware(tokenRepo))

	for _, build := range moduleBuilders {
		build(db).Register(protected)
	}

	return r
}
