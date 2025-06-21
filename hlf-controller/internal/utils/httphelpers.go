package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// ReadJSON reads and unmarshals JSON from the request body
func ReadJSON(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	return json.Unmarshal(body, v)
}

// WriteJSON writes a JSON response with the given status code
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// WriteError writes an error response as JSON
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"status": "error", "message": message})
}
