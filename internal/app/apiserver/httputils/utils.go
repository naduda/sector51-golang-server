package httputils

import (
	"encoding/json"
	"net/http"
)

// Error ...
func SendError(w http.ResponseWriter, code int, err error) {
	Respond(w, code, map[string]string{"error": err.Error()})
}

// Respond ...
func Respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
