PROTO_DIR := protos
GO_OUT := golang
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

.PHONY: proto proto-tools

# Generate Go code from protos into ./golang
proto:
	protoc -I=$(PROTO_DIR) $(PROTO_FILES) \
		--go_out=$(GO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT) --go-grpc_opt=paths=source_relative

# Install required protoc plugins (run once, adds binaries to GOPATH/bin)
proto-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
