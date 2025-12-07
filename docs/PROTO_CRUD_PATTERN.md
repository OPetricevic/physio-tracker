# Proto-First CRUD Pattern (reference)

- **Source of truth**: Define request/response messages in `.proto`. Use `protojson.Unmarshal/Marshal` in REST handlers to avoid ad-hoc structs.
- **Handlers**:
  - Unmarshal JSON into the generated proto request (e.g., `CreatePatientRequest`).
  - Validate required fields manually (proto3 has no `required`).
  - Map to service input (or pass proto if the service uses proto types).
  - Marshal proto responses with `protojson` and `EmitUnpopulated` if needed.
- **Services**:
  - Accept either proto messages or mapped structs; validation/UUID/timestamps live here.
  - Keep business rules and call outbound ports.
- **Outbound ports/repos**:
  - Prefer using proto-generated types (or structs tightly aligned to them) so shapes don’t drift.
  - Implement DB adapters against these ports; swap memory/Postgres without changing service/handler.
- **Routing**:
  - Single HTTP server (mux/chi) with multiple routes; no separate listeners per endpoint.
  - Use query params for search/filter; keep paging in service/repo.

Use this as a reminder to stay proto-first and avoid duplicate request structs.***
