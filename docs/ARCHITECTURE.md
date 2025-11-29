# Backend Layout Conventions (Ports/Adapters)

This documents the naming and folder conventions to mirror the ports/adapters pattern you described (similar to the pickleball repo). Use this as the single reference when adding new domains.

## High-level folders
- `internal/`
  - `services/<Domain>/<Domain>Service.go` — business logic (use-cases) and inbound port for that domain. Tests next to it (`<Domain>Service_Test.go`).
  - `database/<Domain>/<Domain>Repository.go` — DB adapter implementing outbound port. Helper/test files live here too.
  - `api/rest/core/` — transport layer (core API). Inside:
    - `inbound/<Domain>/<Domain>_inbound_port.go` — inbound interfaces the transport calls.
    - `outbound/<Domain>/<Domain>_outbound_port.go` — outbound interfaces the core layer depends on (DB, PDF, backup, etc.).
    - Handlers/controllers for HTTP/REST can sit alongside or under `handlers/` in this area.
- `protos/` — proto definitions (source of truth).
- `golang/` — generated code from protos (gitignored except for `.gitkeep`).

## Naming rules
- One interface per file where practical:
  - Inbound: `<Domain>_inbound_port.go` (methods the handler calls).
  - Outbound: `<Domain>_outbound_port.go` (DB/PDF/Backup/etc.), or more specific names if clearer.
  - Service: `<Domain>Service.go` (service struct + service interface).
- Handlers/transport adapters go under `api/rest/core/...`, implementing the inbound port (or generated gRPC if added later) and delegating to service.

## Example (Patients)
- `internal/services/Patients/PatientsService.go`
- `internal/database/Patients/PatientsRepository.go`
- `internal/api/rest/core/inbound/Patients/Patients_inbound_port.go`
- `internal/api/rest/core/outbound/Patients/Patients_outbound_port.go` (or more specific: `Anamnesis_outbound_port.go`, `PDF_outbound_port.go`, `Backup_outbound_port.go`)
- HTTP handlers can live under `internal/api/rest/core/handlers/Patients/*.go` and call the inbound port.

## ORM note
If using an ORM (e.g., GORM), keep the ORM models/adapters in `internal/database/<Domain>/` implementing the outbound repository ports. Keep domain models (proto-aligned) in the service layer; do not leak ORM models upward.
