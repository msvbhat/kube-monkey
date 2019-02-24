package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Handler for health check endpoint
func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "healthy")
}

// Start a Health Check REST Endpoint for checking health
func healthCheck(status chan string) {
	defer func() { status <- "stop" }()
	router := mux.NewRouter()
	router.HandleFunc("/", health)
	router.HandleFunc("/health", health)
	log.Fatal(http.ListenAndServe(":8080", router))
}
