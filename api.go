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

// TODO: Define and provide meaningful metrics via below Handler
// Handler for metrics endpoint.
func metrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Still not implemented")
}

// Start a Health Check REST Endpoint for checking health
func healthCheck(status chan string) {
	defer func() { status <- "stop" }()
	router := mux.NewRouter()
	router.HandleFunc("/", health)
	router.HandleFunc("/health", health)
	router.HandleFunc("/metrics", metrics)
	log.Fatal(http.ListenAndServe(":8080", router))
}
