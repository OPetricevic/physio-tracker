# Protos

- Source of truth for API/domain messages. We generate Go types into `../golang`.
- Use `protoc` (or buf) locally; generated code is not committed.

Example generation (requires `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`):

```bash
cd backend/protos
protoc --go_out=../golang --go_opt=paths=source_relative \
  --go-grpc_out=../golang --go-grpc_opt=paths=source_relative \
  patients.proto
```

`go_package` is set to `github.com/OPetricevic/physio-tracker/backend/golang/patients`.
