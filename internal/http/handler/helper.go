package handler

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, code int, message string, details ...any) {
	errResp := ErrorResponse{
		Message: message,
		Details: details,
	}

	if len(details) == 1 {
		errResp.Details = details[0]
	}

	writeJson(w, code, errResp)
}

func writeJSONOk(w http.ResponseWriter, data any) {
	writeJson(w, http.StatusOK, data)
}
