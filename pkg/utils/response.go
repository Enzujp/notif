package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	if statusCode > 499 {
		log.Println("Responding with 5XX error: ", message)
	}
	type errorResponse struct {
		Error string `json: "error`
	}

	respondWithJson(w, statusCode, errorResponse{
		Error: message,
	})
}

func respondWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload %v:", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
