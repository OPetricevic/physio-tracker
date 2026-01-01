package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"net/http"

	authh "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers/auth"
	authctrl "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/inbound/auth"
	"github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/middleware"
	mwauth "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/middleware"
	dbauth "github.com/OPetricevic/physio-tracker/backend/internal/database/auth"
	dbcreds "github.com/OPetricevic/physio-tracker/backend/internal/database/doctors"
	dbdoctors "github.com/OPetricevic/physio-tracker/backend/internal/database/doctors"
	authsvc "github.com/OPetricevic/physio-tracker/backend/internal/services/auth"
	"gorm.io/gorm"
)

// BuildRouter constructs the HTTP router with health, auth (public), and protected routes.
func BuildRouter(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()
	// API subrouter under /api
	api := r.PathPrefix("/api").Subrouter()
	api.Use(mwauth.RecoverMiddleware)
	api.Use(mwauth.LoggingMiddleware)
	api.Use(middleware.CORSMiddleware)

	// Health
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
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
	authHandler.RegisterRoutes(api)

	// Serve uploaded branding assets
	staticDir := http.Dir("uploads")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Protected routes: everything else goes under a subrouter with auth middleware.
	protected := api.PathPrefix("/").Subrouter()
	protected.Use(mwauth.AuthMiddleware(tokenRepo))

	for _, build := range moduleBuilders {
		build(db).Register(protected)
	}

	// Serve frontend build (if present)
	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		frontendDir = filepath.Join("frontend", "dist")
	}
	spa := spaHandler(frontendDir)
	r.PathPrefix("/").Handler(spa)

	return r
}

// spaHandler serves static files and falls back to index.html for SPA routes.
func spaHandler(dir string) http.Handler {
	fileServer := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If file exists, serve it
		path := filepath.Join(dir, r.URL.Path)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}
		// fallback to index.html
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}
