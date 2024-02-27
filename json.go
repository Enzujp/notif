package main

import (
	// "encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Println("Responding with 5XX error: ", message)
	}
}

// func respondWithJson()
