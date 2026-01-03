package main

import (
	"log"
	"net/http"
	"os"
	"time"

	corehandlers "github.com/OPetricevic/physio-tracker/backend/internal/api/rest/core/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	addr := envOrDefault("PORT", "3600")

	dsn := envOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/physio?sslmode=disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Keep GORM output quiet (only errors); we log requests separately.
		Logger: logger.Default.LogMode(logger.Error),
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// Build router (health, auth, protected modules).
	r := corehandlers.BuildRouter(db)

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
