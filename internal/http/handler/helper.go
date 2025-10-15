package handler

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data) // nolint: errcheck
}

func Error(w http.ResponseWriter, code int, message string, details ...any) {
	errResp := ErrorResponse{
		Message: message,
		Details: details,
	}

	if len(details) == 1 {
		errResp.Details = details[0]
	}

	writeJSON(w, code, errResp)
}

func Ok(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, data)
}
