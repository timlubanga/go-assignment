package main

import (
	"log"
	"net/http"

	"time"

	"go-api-assignment/handler"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.Path("/api/helm-images").Methods(http.MethodPost).HandlerFunc(handler.HelmImagesHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(handler.MethodNotAllowedHandler)


	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 1500 * time.Second,
		ReadTimeout:  1500 * time.Second,
	}

	log.Printf("Server starting on address: %v", srv.Addr)

	log.Fatal(srv.ListenAndServe())

}
