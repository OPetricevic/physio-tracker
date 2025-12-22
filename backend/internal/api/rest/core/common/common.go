package common

import (
	"encoding/json"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// local protojson helper
var JSONPB = &protojson.UnmarshalOptions{DiscardUnknown: true}

func WriteProto(w http.ResponseWriter, msg proto.Message, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// UseProtoNames ensures we keep proto snake_case field names to match FE DTOs.
	b, err := protojson.MarshalOptions{EmitUnpopulated: true, UseEnumNumbers: true, UseProtoNames: true}.Marshal(msg)
	if err != nil {
		http.Error(w, "internal_error: failed to encode response", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(b)
}

func WriteJSONError(w http.ResponseWriter, code, message string, status int) {
	WriteJSON(w, map[string]string{
		"error":   code,
		"message": message,
	}, status)
}

func WriteJSON(w http.ResponseWriter, payload interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "internal_error: failed to encode response", http.StatusInternalServerError)
	}
}

// ValidateProto calls generated PGV validators if available.
func ValidateProto(msg interface{}) error {
	type validator interface{ ValidateAll() error }
	if v, ok := msg.(validator); ok {
		return v.ValidateAll()
	}
	type validatorSimple interface{ Validate() error }
	if v, ok := msg.(validatorSimple); ok {
		return v.Validate()
	}
	return nil
}
