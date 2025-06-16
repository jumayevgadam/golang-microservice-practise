package reqvalidator

import (
	"encoding/json"
	"log"
	"net/http"
)

func Respond(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("json.NewEncoder.Encode: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	Respond(w, status, map[string]string{"error": message})
}
