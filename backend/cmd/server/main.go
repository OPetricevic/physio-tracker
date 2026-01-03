package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
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

	port := choosePort(addr, 15)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	l, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("failed to bind to %s: %v", server.Addr, err)
	}

	log.Printf("listening on %s", server.Addr)
	if err := server.Serve(l); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// choosePort tries the requested port first, then the next N-1 ports, and returns the first free one.
func choosePort(start string, attempts int) string {
	port := start
	startNum, err := strconv.Atoi(start)
	if err != nil {
		// not a number, just return as-is
		return start
	}
	for i := 0; i < attempts; i++ {
		p := startNum + i
		addr := fmt.Sprintf(":%d", p)
		l, err := net.Listen("tcp", addr)
		if err == nil {
			_ = l.Close()
			return strconv.Itoa(p)
		}
	}
	// fallback to original even if busy; Serve will fail loudly
	return start
}
