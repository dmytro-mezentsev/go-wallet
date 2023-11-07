package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"wallet.com/wallet/wallet/internal/api/middleware"
	"wallet.com/wallet/wallet/internal/api/routes"
)

func main() {

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	routes.HelloWorldRoute(r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	//keep alive
	log.Fatal(srv.ListenAndServe())
}
