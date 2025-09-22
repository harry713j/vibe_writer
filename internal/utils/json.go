package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON[T any](w http.ResponseWriter, code int, payload T) {
	w.Header().Add("Content-Type", "application/json")

	data, err := json.Marshal(payload)

	if err != nil {
		log.Println("Failed to marshal json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {

		log.Printf("Error writing response: %v\n", err)
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Println("Server error: ", message)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errorResponse{
		Error: message,
	})

}
