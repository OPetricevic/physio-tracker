FRONTEND_DIR := frontend
BACKEND_DIR := backend
# default app creds (created by bootstrap script)
DB_URL ?= postgres://physio:physio@localhost:5433/physio?sslmode=disable
RELEASE_DIR := release/physio-bundle
RELEASE_DIR_WIN := release/physio-bundle-win
RELEASE_DIR_WIN_X64 := release/physio-bundle-win-x64
RELEASE_DIR_WIN_ARM64 := release/physio-bundle-win-arm64
PROTO_DIR := $(BACKEND_DIR)/protos
PROTO_OUT := $(BACKEND_DIR)/golang
PROTO_OUT_PKG := $(PROTO_OUT)/github.com/OPetricevic/physio-tracker/backend/golang/patients
PROTOC := protoc
PROTOC_GEN_GO := $(shell go env GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GORM := $(shell go env GOPATH)/bin/protoc-gen-gorm
PROTOC_GEN_VALIDATE := $(shell go env GOPATH)/bin/protoc-gen-validate
# Fixed include path for PGV (pinned to v0.6.13)
VALIDATE_INC := $(shell go env GOPATH)/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v0.6.13

.PHONY: frontend-install frontend-dev frontend-build run backend-run backend-migrate backend-bootstrap backend-proto

frontend-install:
	npm --prefix $(FRONTEND_DIR) install

frontend-dev:
	npm --prefix $(FRONTEND_DIR) run dev

frontend-build:
	npm --prefix $(FRONTEND_DIR) run build

# Convenience alias for local development
run: frontend-dev

# Backend: run the server (expects DB already created/migrated)
backend-run:
	cd $(BACKEND_DIR) && DATABASE_URL=$(DB_URL) PORT=3600 go run ./cmd/server

# Backend: build binary
backend-build:
	cd $(BACKEND_DIR) && DATABASE_URL=$(DB_URL) PORT=3600 go build -o ../$(RELEASE_DIR)/server ./cmd/server

# Backend: apply SQL migrations in order using psql (requires DB to exist)
backend-migrate:
	cd $(BACKEND_DIR) && DATABASE_URL=$(DB_URL) ./scripts/migrate.sh

# Backend: create DB/user (if missing) using the default superuser connection
# WARNING: adjust DB_BOOTSTRAP_URL/DB_NAME/DB_APP_USER/DB_APP_PASS to your environment.
backend-bootstrap:
	cd $(BACKEND_DIR) && ./scripts/bootstrap_postgres.sh

# Backend: generate Go code from protos
backend-proto:
	@command -v $(PROTOC) >/dev/null || { echo "protoc not found; install protoc"; exit 1; }
	@test -x $(PROTOC_GEN_GO) || { echo "protoc-gen-go missing; run: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"; exit 1; }
	@test -x $(PROTOC_GEN_GORM) || { echo "protoc-gen-gorm missing; run: go install github.com/infobloxopen/protoc-gen-gorm@latest"; exit 1; }
	@test -x $(PROTOC_GEN_VALIDATE) || { echo "protoc-gen-validate missing; run: go install github.com/envoyproxy/protoc-gen-validate/cmd/protoc-gen-validate@latest"; exit 1; }
	@test -d "$(VALIDATE_INC)" || { echo "validate proto include not found; ensure protoc-gen-validate v0.6.13 is downloaded"; exit 1; }
	$(PROTOC) -I $(PROTO_DIR) -I $(VALIDATE_INC) --go_out=$(PROTO_OUT) --gorm_out=$(PROTO_OUT) --validate_out="lang=go,paths=source_relative:$(PROTO_OUT_PKG)" $(PROTO_DIR)/*.proto

# Run backend and frontend together (expects DB ready and npm install done)
dev:
	./scripts/dev.sh

# Alias: generate protos (Go + gorm)
proto: backend-proto

# Package: build frontend + Windows x64 bundle + zip for release
package: clean-release frontend-build
	@mkdir -p $(RELEASE_DIR)
	# build Windows x64 backend binary into release dir
	cd $(BACKEND_DIR) && GOOS=windows GOARCH=amd64 DATABASE_URL=$(DB_URL) PORT=3600 go build -o ../$(RELEASE_DIR)/server.exe ./cmd/server
	# copy frontend build
	mkdir -p $(RELEASE_DIR)/frontend
	cp -r frontend/dist $(RELEASE_DIR)/frontend/
	# copy assets (fonts), uploads placeholder, scripts, migrations
	mkdir -p $(RELEASE_DIR)/assets/fonts
	cp -r backend/assets/fonts/* $(RELEASE_DIR)/assets/fonts/
	mkdir -p $(RELEASE_DIR)/uploads
	mkdir -p $(RELEASE_DIR)/migrations
	cp backend/migrations/*.sql $(RELEASE_DIR)/migrations/
	cp scripts/start_windows.ps1 $(RELEASE_DIR)/
	mkdir -p $(RELEASE_DIR)/scripts/win
	cp scripts/win/*.ps1 $(RELEASE_DIR)/scripts/win/ 2>/dev/null || true
	cp scripts/win/*.iss $(RELEASE_DIR)/scripts/win/ 2>/dev/null || true
	@echo "Windows bundle created at $(RELEASE_DIR)"
	# zip the bundle for release upload
	cd release && zip -r physio-windows-x64.zip physio-bundle
	@echo "Zipped bundle at release/physio-windows-x64.zip"

clean-release:
	rm -rf release

# Cross-compile Windows portable bundle and zip it
package-win: clean-release frontend-build
	@mkdir -p $(RELEASE_DIR_WIN)
	# cross-compile backend to Windows
	cd $(BACKEND_DIR) && GOOS=windows GOARCH=amd64 DATABASE_URL=$(DB_URL) PORT=3600 go build -o ../$(RELEASE_DIR_WIN)/server.exe ./cmd/server
	# copy frontend build
	mkdir -p $(RELEASE_DIR_WIN)/frontend
	cp -r frontend/dist $(RELEASE_DIR_WIN)/frontend/
	# copy assets, migrations, uploads placeholder, scripts
	mkdir -p $(RELEASE_DIR_WIN)/assets/fonts
	cp -r backend/assets/fonts/* $(RELEASE_DIR_WIN)/assets/fonts/
	mkdir -p $(RELEASE_DIR_WIN)/uploads
	mkdir -p $(RELEASE_DIR_WIN)/migrations
	cp backend/migrations/*.sql $(RELEASE_DIR_WIN)/migrations/
	cp scripts/start_windows.ps1 $(RELEASE_DIR_WIN)/
	mkdir -p $(RELEASE_DIR_WIN)/scripts/win
	cp scripts/win/*.ps1 $(RELEASE_DIR_WIN)/scripts/win/ 2>/dev/null || true
	# zip the bundle
	cd release && zip -r physio-windows-portable.zip physio-bundle-win
	@echo "Windows portable bundle: release/physio-windows-portable.zip"

# Build Windows bundle for x64 installer (server.exe only, no zip)
package-win-x64: clean-release frontend-build
	@mkdir -p $(RELEASE_DIR_WIN_X64)
	cd $(BACKEND_DIR) && GOOS=windows GOARCH=amd64 DATABASE_URL=$(DB_URL) PORT=3600 go build -o ../$(RELEASE_DIR_WIN_X64)/server.exe ./cmd/server
	mkdir -p $(RELEASE_DIR_WIN_X64)/frontend
	cp -r frontend/dist $(RELEASE_DIR_WIN_X64)/frontend/
	mkdir -p $(RELEASE_DIR_WIN_X64)/assets/fonts
	cp -r backend/assets/fonts/* $(RELEASE_DIR_WIN_X64)/assets/fonts/
	mkdir -p $(RELEASE_DIR_WIN_X64)/uploads
	mkdir -p $(RELEASE_DIR_WIN_X64)/migrations
	cp backend/migrations/*.sql $(RELEASE_DIR_WIN_X64)/migrations/
	cp scripts/start_windows.ps1 $(RELEASE_DIR_WIN_X64)/
	mkdir -p $(RELEASE_DIR_WIN_X64)/scripts/win
	cp scripts/win/*.ps1 $(RELEASE_DIR_WIN_X64)/scripts/win/ 2>/dev/null || true
	cp scripts/win/*.iss $(RELEASE_DIR_WIN_X64)/scripts/win/ 2>/dev/null || true
	@echo "Windows x64 bundle: $(RELEASE_DIR_WIN_X64)"

# Build Windows bundle for ARM64 installer (server.exe only, no zip)
package-win-arm64: clean-release frontend-build
	@mkdir -p $(RELEASE_DIR_WIN_ARM64)
	cd $(BACKEND_DIR) && GOOS=windows GOARCH=arm64 DATABASE_URL=$(DB_URL) PORT=3600 go build -o ../$(RELEASE_DIR_WIN_ARM64)/server.exe ./cmd/server
	mkdir -p $(RELEASE_DIR_WIN_ARM64)/frontend
	cp -r frontend/dist $(RELEASE_DIR_WIN_ARM64)/frontend/
	mkdir -p $(RELEASE_DIR_WIN_ARM64)/assets/fonts
	cp -r backend/assets/fonts/* $(RELEASE_DIR_WIN_ARM64)/assets/fonts/
	mkdir -p $(RELEASE_DIR_WIN_ARM64)/uploads
	mkdir -p $(RELEASE_DIR_WIN_ARM64)/migrations
	cp backend/migrations/*.sql $(RELEASE_DIR_WIN_ARM64)/migrations/
	cp scripts/start_windows.ps1 $(RELEASE_DIR_WIN_ARM64)/
	mkdir -p $(RELEASE_DIR_WIN_ARM64)/scripts/win
	cp scripts/win/*.ps1 $(RELEASE_DIR_WIN_ARM64)/scripts/win/ 2>/dev/null || true
	cp scripts/win/*.iss $(RELEASE_DIR_WIN_ARM64)/scripts/win/ 2>/dev/null || true
	@echo "Windows ARM64 bundle: $(RELEASE_DIR_WIN_ARM64)"
