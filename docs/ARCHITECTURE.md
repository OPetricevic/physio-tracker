# Backend Layout Conventions (Ports/Adapters)

This documents the naming and folder conventions to mirror the ports/adapters pattern you described (similar to the pickleball repo). Use this as the single reference when adding new domains.

## High-level folders
- `internal/`
  - `services/<Domain>/<Domain>Service.go` — business logic (use-cases) and inbound port interface for that domain.
  - `ports/inbound/<Domain>/<Domain>InboundPort.go` — interface(s) the transport calls (mirrors service methods if you separate them).
  - `ports/outbound/<Domain>/<Domain>Repository.go` (one file per outbound interface, e.g., PatientRepository, AnamnesisRepository, PDFGenerator, BackupStore).
  - `database/<Domain>/<Domain>Repository.go` (DB adapter implementing outbound port; helpers/tests live alongside, e.g., `UsersRepositoryHelpers.go`, `UsersRepositoryTest.go`).
  - `transport/http/<Domain>/...` or `transport/grpc/<Domain>/...` — handlers/adapters implementing inbound ports (files like `PatientsInbound.go` or `UsersInbound.go`).
- `protos/` — proto definitions (source of truth).
- `golang/` — generated code from protos (gitignored except for `.gitkeep`).

## Naming rules
- Keep one interface per file where practical:
  - Inbound: `<Domain>InboundPort.go` (methods the handler calls).
  - Outbound: `<Domain>Repository.go`, `<Domain>PDFGenerator.go`, etc.
  - Service: `<Domain>Service.go` (defines service struct and the service interface if kept here).
- Handlers/transport adapters should live under `transport/<protocol>/<Domain>/` and implement the inbound port (or generated gRPC interfaces), delegating to service/controller.
- Controllers (if used) can sit next to transport adapters: e.g., `transport/http/<Domain>/<Domain>Controller.go`.

## Example (Patients)
- `internal/services/Patients/PatientsService.go` — business logic + PatientsService interface.
- `internal/ports/inbound/Patients/PatientsInboundPort.go` — inbound interface called by HTTP/gRPC.
- `internal/ports/outbound/Patients/PatientsRepository.go` — outbound interface for DB.
- `internal/ports/outbound/Patients/AnamnesisRepository.go`
- `internal/ports/outbound/Patients/PDFGenerator.go`
- `internal/ports/outbound/Patients/BackupStore.go`
- `internal/database/Patients/PatientsRepository.go` — Postgres (ORM) implementation.
- `internal/transport/http/Patients/PatientsInbound.go` — HTTP handler implementing inbound port.
- `internal/transport/grpc/Patients/PatientsInbound.go` — gRPC server implementing inbound port/generated interface.

## ORM note
If using an ORM (e.g., GORM), keep the ORM models/adapters in `internal/database/<Domain>/` implementing the outbound repository ports. Keep domain models (proto-aligned) in the service layer; do not leak ORM models upward.
